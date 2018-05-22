package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/sanbornm/go-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

type HTTPRequester struct {
}

// Fetch will return an HTTP request to the specified url and return
// the body of the result. An error will occur for a non 200 status code.
func (httpRequester *HTTPRequester) Fetch(url string) (io.ReadCloser, error) {
	fmt.Printf("GET %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad http status from %s: %v", url, resp.Status)
	}

	return resp.Body, nil
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the ranch CLI binary",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		updater := &selfupdate.Updater{
			CurrentVersion: Version,
			ApiURL:         "http://ranch-updates.goodeggs.com/stable/",
			BinURL:         "http://ranch-updates.goodeggs.com/stable/",
			DiffURL:        "http://ranch-updates.goodeggs.com/stable/",
			Dir:            ".ranch-selfupdate/",
			CmdName:        "ranch",
			Requester:      &HTTPRequester{},
			ForceCheck:     true,
		}

		return updater.BackgroundRun()
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
