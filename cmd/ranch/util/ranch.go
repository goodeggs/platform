package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/parnurzeal/gorequest"
	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/spf13/viper"
)

type RanchConfig struct {
	Name      string                `json:"name"`
	EnvId     string                `json:"env_id"`
	Processes RanchConfigProcessMap `json:"processes"`
}

type RanchConfigProcess struct {
	Command   string `json:"command"`
	Count     int    `json:"count"`
	Instances int    `json:"instances"` // deprecated
	Memory    int    `json:"memory"`
}

type RanchConfigProcessMap map[string]RanchConfigProcess

type RanchFormationEntry struct {
	Balancer string `json:"balancer"`
	Count    int    `json:"count"`
	Memory   int    `json:"memory"`
}

type RanchFormation map[string]RanchFormationEntry

type Process struct {
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

type Processes []Process

type Release struct {
	Id      string    `json:"id"`
	App     string    `json:"app"`
	Created time.Time `json:"created"`
	Status  string    `json:"status"`
}

type Releases []Release

var ValidAppName = regexp.MustCompile(`\A[a-z][-a-z0-9]{3,29}\z`)
var ValidProcessName = regexp.MustCompile(`\A[a-z][-a-z0-9]{2,29}\z`)

func getClient(authToken string) *gorequest.SuperAgent {
	return jsonClient().
		SetBasicAuth(authToken, "x-auth-token")
}

func RanchValidateConfig(config *RanchConfig) (errors []error) {
	if !ValidAppName.MatchString(config.Name) {
		errors = append(errors, fmt.Errorf("app name '%s' is invalid: must match %s", config.Name, ValidAppName.String()))
	}

	for name, _ := range config.Processes {
		if !ValidProcessName.MatchString(name) {
			errors = append(errors, fmt.Errorf("process name '%s' is invalid: must match %s", name, ValidProcessName.String()))
		}
	}

	return errors
}

func RanchLoadSettings() (err error) {
	authToken := viper.GetString("token")
	url := fmt.Sprintf("%s/settings", viper.Get("endpoint"))

	resp, body, errs := getClient(authToken).Get(url).End()

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
