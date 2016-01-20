package util

import (
	"encoding/json"
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type EcruRelease struct {
	Id            string `json:"id"`
	ProjectId     string `json:"project"`
	ConvoxRelease string `json:"convoxRelease"`
	Sha           string `json:"sha"`
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
		return nil, fmt.Errorf("Error fetching releases from Ecru: status code %s", resp.StatusCode)
	}

	var ecruReleases []EcruRelease
	err := json.Unmarshal([]byte(body), &ecruReleases)

	if err != nil {
		return nil, err
	}

	return ecruReleases, nil
}
