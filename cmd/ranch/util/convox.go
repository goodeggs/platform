package util

import (
	"fmt"
	"os"

	"github.com/convox/rack/client"
	"github.com/spf13/viper"
)

func Convox() *client.Client {
	if !viper.IsSet("convox.host") || !viper.IsSet("convox.password") {
		fmt.Println("must set 'convox.host' and 'convox.password' in $HOME/.ranch.yaml")
		os.Exit(1)
	}

	host := viper.GetString("convox.host")
	password := viper.GetString("convox.password")
	version := "20151211151200"
	return client.New(host, password, version)
}
