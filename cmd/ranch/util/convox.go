package util

import (
	"fmt"

	"github.com/convox/rack/client"
	"github.com/spf13/viper"
)

func convoxClient() (*client.Client, error) {
	if !viper.IsSet("convox.host") || !viper.IsSet("convox.password") {
		return nil, fmt.Errorf("must set 'convox.host' and 'convox.password' in $HOME/.ranch.yaml")
	}

	host := viper.GetString("convox.host")
	password := viper.GetString("convox.password")
	version := "20151211151200"
	return client.New(host, password, version), nil
}

func ConvoxPs(appName string) (Processes, error) {
	client, err := convoxClient()

	if err != nil {
		return nil, err
	}

	convoxPs, err := client.GetProcesses(appName, false) // false == no process stats

	if err != nil {
		return nil, err
	}

	var ps Processes

	for _, v := range convoxPs {
		p := Process(v)
		ps = append(ps, p)
	}

	return ps, nil
}
