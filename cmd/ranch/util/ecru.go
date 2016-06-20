package util

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/parnurzeal/gorequest"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/viper"
)

type EcruError struct {
	Message string `json:"error"`
}

type EcruRelease struct {
	Id            string `json:"_id"`
	ProjectId     string `json:"project"`
	ConvoxRelease string `json:"convoxRelease"`
	Sha           string `json:"sha"`
}

type EcruSecret struct {
	Id string `json:"_id"`
}

func ecruUrl(pathname string) string {
	u, _ := url.Parse("https://ecru.goodeggs.com/api/v1")
	u.Path = path.Join(u.Path, pathname)
	return u.String()
}

func ecruClient() (*gorequest.SuperAgent, error) {
	if !viper.IsSet("convox.password") {
		return nil, fmt.Errorf("must set 'convox.password' in $HOME/.ranch.yaml")
	}

	request := jsonClient().
		SetBasicAuth(viper.GetString("convox.password"), "x-auth-token")

	return request, nil
}

func EcruCreateRelease(appName, sha, convoxRelease string) error {

	client, err := ecruClient()

	if err != nil {
		return err
	}

	pathname := fmt.Sprintf("/projects/%s/releases", appName)
	reqBody := fmt.Sprintf(`{"sha":"%s","convoxRelease":"%s"}`, sha, convoxRelease)

	resp, body, errs := client.Post(ecruUrl(pathname)).Send(reqBody).End()

	if len(errs) > 0 {
		return errs[0]
	}

	makeError := func(statusCode int, message string) error {
		return fmt.Errorf("Error creating Ecru release [HTTP %d]: %s", statusCode, message)
	}

	switch resp.StatusCode {
	case 201:
		return nil
	case 400:
		var ecruError EcruError
		err := json.Unmarshal([]byte(body), &ecruError)
		if err == nil {
			return makeError(resp.StatusCode, ecruError.Message)
		}
	}

	return makeError(resp.StatusCode, body)
}

func EcruReleases(appName string) ([]EcruRelease, error) {

	client, err := ecruClient()

	if err != nil {
		return nil, err
	}

	pathname := fmt.Sprintf("/projects/%s/releases", appName)

	resp, body, errs := client.Get(ecruUrl(pathname)).End()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error fetching releases from Ecru: status code %d", resp.StatusCode)
	}

	var ecruReleases []EcruRelease
	err = json.Unmarshal([]byte(body), &ecruReleases)

	if err != nil {
		return nil, err
	}

	return ecruReleases, nil
}

func EcruGetSecret(appName, secretId string) (plaintext string, err error) {

	client, err := ecruClient()

	if err != nil {
		return "", err
	}

	pathname := fmt.Sprintf("/projects/%s/secrets/%s", appName, secretId)

	resp, body, errs := client.Get(ecruUrl(pathname)).End()

	if len(errs) > 0 {
		return "", errs[0]
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error fetching secret from Ecru: status code %d", resp.StatusCode)
	}

	return body, nil
}

func EcruCreateSecret(appName, plaintext string) (secretId string, err error) {

	client, err := ecruClient()

	if err != nil {
		return "", err
	}

	pathname := fmt.Sprintf("/projects/%s/secrets", appName)

	resp, body, errs := client.
		Post(ecruUrl(pathname)).
		Set("Content-Type", "text/plain").
		Send(plaintext).
		End()

	if len(errs) > 0 {
		return "", errs[0]
	}

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("Error creating secret in Ecru: status code %d", resp.StatusCode)
	}

	var ecruSecret EcruSecret
	err = json.Unmarshal([]byte(body), &ecruSecret)

	if err != nil {
		return "", err
	}

	return ecruSecret.Id, nil
}
