package util

import (
	"fmt"
	"log"
	"time"

	"github.com/goodeggs/platform/cmd/ranch/Godeps/_workspace/src/github.com/parnurzeal/gorequest"
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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
