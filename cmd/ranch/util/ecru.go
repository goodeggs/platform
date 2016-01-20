package util

import (
	"encoding/json"
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type EcruError struct {
	Message string `json:"error"`
}

type EcruRelease struct {
	Id            string `json:"id"`
	ProjectId     string `json:"project"`
	ConvoxRelease string `json:"convoxRelease"`
	Sha           string `json:"sha"`
}

func EcruCreateRelease(appName, sha, convoxRelease string) error {

	request := gorequest.New()

	url := fmt.Sprintf("https://ecru.goodeggs.com/api/v1/projects/%s/releases", appName)

	noRedirects := func(req gorequest.Request, via []gorequest.Request) error {
		return fmt.Errorf("refusing to follow redirect")
	}

	resp, body, errs := request.
		Post(url).
		RedirectPolicy(noRedirects).
		Send(fmt.Sprintf(`{"sha":"%s","convoxRelease":"%s"}`, sha, convoxRelease)).
		End()

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

	request := gorequest.New()

	url := fmt.Sprintf("https://ecru.goodeggs.com/api/v1/projects/%s/releases", appName)

	noRedirects := func(req gorequest.Request, via []gorequest.Request) error {
		return fmt.Errorf("refusing to follow redirect")
	}

	resp, body, errs := request.
		Get(url).
		RedirectPolicy(noRedirects).
		End()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error fetching releases from Ecru: status code %d", resp.StatusCode)
	}

	var ecruReleases []EcruRelease
	err := json.Unmarshal([]byte(body), &ecruReleases)

	if err != nil {
		return nil, err
	}

	return ecruReleases, nil
}
