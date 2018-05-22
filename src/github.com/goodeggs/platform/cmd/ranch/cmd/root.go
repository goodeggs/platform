package cmd

import (
	"fmt"
	"os"

	"github.com/goodeggs/platform/cmd/ranch/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var App string
var Version string

var RootCmd = &cobra.Command{
	Use:   "ranch",
	Short: "Ranch CLI",
	Long: `A CLI interface to Ranch aka the Good Eggs platform,
  maintained with love by the Delivery Engineering team in Go.
  More information is available at https://github.com/goodeggs/platform`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVarP(&util.Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVarP(&App, "app", "a", "", "application name")
	RootCmd.PersistentFlags().StringP("filename", "f", "", "config filename (defaults to .ranch.yaml)")
	RootCmd.SilenceUsage = true
	RootCmd.SilenceErrors = true
}

func initConfig() {
	viper.SetEnvPrefix("ranch")
	viper.SetDefault("endpoint", "https://ranch-api.goodeggs.com")
	viper.BindEnv("endpoint")
	viper.BindEnv("token")

	if err := util.RanchLoadSettings(); err != nil {
		fmt.Println("error trying to authenticate with ranch - did you set RANCH_TOKEN?")
		fmt.Println(err)
		os.Exit(1)
	}
}
