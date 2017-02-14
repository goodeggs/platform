package util

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

func noRedirects(req gorequest.Request, via []gorequest.Request) error {
	return fmt.Errorf("refusing to follow redirect")
}

func jsonClient() *gorequest.SuperAgent {
	return gorequest.New().
		RedirectPolicy(noRedirects).
		Set("Accept", "application/json").
		Set("Content-Type", "application/json")
}

var Verbose bool
