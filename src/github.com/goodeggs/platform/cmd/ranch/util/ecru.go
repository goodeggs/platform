package util

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
)

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
