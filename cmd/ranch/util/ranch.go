package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"regexp"
	"time"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/parnurzeal/gorequest"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/viper"
)

type RanchApiError struct {
	Message string `json:"error"`
}

type RanchApiRelease struct {
	Id            string `json:"id"` // sha
	App           string `json:"app"`
	ConvoxRelease string `json:"convoxRelease"`
}

type RanchApiSecret struct {
	Id      string `json:"_id"`
	Content string `json:"content"`
}

type RanchConfig struct {
	AppName   string                        `json:"name"`
	ImageName string                        `json:"image_name"`
	EnvId     string                        `json:"env_id"`
	Processes map[string]RanchConfigProcess `json:"processes"`
}

type RanchConfigProcess struct {
	Command   string `json:"command"`
	Count     int    `json:"count"`
	Instances int    `json:"instances"` // deprecated
	Memory    int    `json:"memory"`
}

type RanchFormationEntry struct {
	Balancer string `json:"balancer"`
	Count    int    `json:"count"`
	Memory   int    `json:"memory"`
}

type RanchFormation map[string]RanchFormationEntry

type RanchProcess struct {
	Id      string    `json:"id"`
	App     string    `json:"app"`
	Command string    `json:"command"`
	Host    string    `json:"host"`
	Image   string    `json:"image"`
	Name    string    `json:"name"`
	Ports   []string  `json:"ports"`
	Release string    `json:"release"`
	Cpu     float64   `json:"cpu"`
	Memory  float64   `json:"memory"`
	Started time.Time `json:"started"`
}

type RanchRelease struct {
	Id      string    `json:"id"`
	App     string    `json:"app"`
	Created time.Time `json:"created"`
	Status  string    `json:"status"`
}

var ValidAppName = regexp.MustCompile(`\A[a-z][-a-z0-9]{3,29}\z`)
var ValidProcessName = regexp.MustCompile(`\A[a-z][-a-z0-9]{2,29}\z`)

func ranchUrl(pathname string) string {
	u, _ := url.Parse(viper.GetString("endpoint"))
	u.Path = path.Join(u.Path, pathname)
	return u.String()
}

func ranchClient() *gorequest.SuperAgent {
	authToken := viper.GetString("token")
	return jsonClient().
		SetBasicAuth(authToken, "x-auth-token")
}

func RanchValidateConfig(config *RanchConfig) (errors []error) {
	if !ValidAppName.MatchString(config.AppName) {
		errors = append(errors, fmt.Errorf("app name '%s' is invalid: must match %s", config.AppName, ValidAppName.String()))
	}

	if !ValidAppName.MatchString(config.ImageName) {
		errors = append(errors, fmt.Errorf("image name '%s' is invalid: must match %s", config.ImageName, ValidAppName.String()))
	}

	for name, _ := range config.Processes {
		if !ValidProcessName.MatchString(name) {
			errors = append(errors, fmt.Errorf("process name '%s' is invalid: must match %s", name, ValidProcessName.String()))
		}
		if name == "run" {
			errors = append(errors, fmt.Errorf("process name 'run' is invalid: 'run' is a reserved process name"))
		}
	}

	return errors
}

func RanchLoadSettings() (err error) {
	resp, body, errs := ranchClient().Get(ranchUrl("/settings")).End()

	if len(errs) > 0 {
		return errs[0]
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected HTTP response [%d]: %s", resp.StatusCode, body)
	}

	viper.SetConfigType("json")
	if err = viper.ReadConfig(bytes.NewBuffer([]byte(body))); err != nil {
		return err
	}

	return // success
}

func RanchUpdateEnvId(ranchFile, envId string) (err error) {
	contents, err := ioutil.ReadFile(ranchFile)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(`(?m)^(\s*env_id\s*:\s*)(['"\w]+)?(.*)$`)
	if err != nil {
		return err
	}

	updatedContents := re.ReplaceAll(contents, []byte("${1}"+envId+"${3}"))
	if bytes.Equal(updatedContents, contents) {
		// if we didn't find it, we'll prepend
		updatedContents = bytes.Join([][]byte{[]byte("env_id: " + envId), contents}, []byte("\n"))
	}

	err = ioutil.WriteFile(ranchFile, updatedContents, 0644)
	if err != nil {
		return err
	}

	return nil
}

func RanchGetSecret(appName, secretId string) (string, error) {

	client := ranchClient()

	pathname := fmt.Sprintf("/v1/apps/%s/secrets/%s", appName, secretId)

	resp, body, errs := client.Get(ranchUrl(pathname)).End()

	if len(errs) > 0 {
		return "", errs[0]
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error fetching secret from ranch-api: status code %d", resp.StatusCode)
	}

	var secret RanchApiSecret
	if err := json.Unmarshal([]byte(body), &secret); err != nil {
		return "", err
	}

	plaintextBytes, err := base64.StdEncoding.DecodeString(secret.Content)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}

func RanchReleaseExists(appName, sha string) (exists bool, err error) {
	client := ranchClient()

	url := fmt.Sprintf("/v1/apps/%s/releases/%s", appName, sha)
	resp, _, errs := client.Get(ranchUrl(url)).End()

	if len(errs) > 0 {
		return false, errs[0]
	} else if resp.StatusCode == 404 {
		return false, nil
	} else if resp.StatusCode == 200 {
		return true, nil
	}

	return false, fmt.Errorf("error fetching release info: HTTP %d", resp.StatusCode)
}

func RanchCreateRelease(appName, sha, convoxRelease string) error {

	client := ranchClient()

	pathname := fmt.Sprintf("/v1/apps/%s/releases", appName)
	reqBody := fmt.Sprintf(`{"id":"%s","convoxRelease":"%s"}`, sha, convoxRelease)

	resp, body, errs := client.Post(ranchUrl(pathname)).Send(reqBody).End()

	if len(errs) > 0 {
		return errs[0]
	}

	makeError := func(statusCode int, message string) error {
		return fmt.Errorf("Error creating Ranch release [HTTP %d]: %s", statusCode, message)
	}

	switch resp.StatusCode {
	case 201:
		return nil
	case 400:
		var ranchError RanchApiError
		err := json.Unmarshal([]byte(body), &ranchError)
		if err == nil {
			return makeError(resp.StatusCode, ranchError.Message)
		}
	}

	return makeError(resp.StatusCode, body)
}

func RanchReleases(appName string) ([]RanchApiRelease, error) {

	client := ranchClient()

	pathname := fmt.Sprintf("/v1/apps/%s/releases", appName)

	resp, body, errs := client.Get(ranchUrl(pathname)).End()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error fetching releases from Ranch: status code %d", resp.StatusCode)
	}

	var ranchReleases []RanchApiRelease
	err := json.Unmarshal([]byte(body), &ranchReleases)

	if err != nil {
		return nil, err
	}

	return ranchReleases, nil
}

func RanchCreateSecret(appName, plaintext string) (secretId string, err error) {

	client := ranchClient()

	pathname := fmt.Sprintf("/v1/apps/%s/secrets", appName)

	secret := RanchApiSecret{
		Content: base64.StdEncoding.EncodeToString([]byte(plaintext)),
	}

	resp, body, errs := client.
		Post(ranchUrl(pathname)).
		Send(secret).
		End()

	if len(errs) > 0 {
		return "", errs[0]
	} else if resp.StatusCode != 201 {
		return "", fmt.Errorf("Error creating secret in ranch-api: status code %d", resp.StatusCode)
	}

	if err = json.Unmarshal([]byte(body), &secret); err != nil {
		return "", err
	}

	return secret.Id, nil
}
