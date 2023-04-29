package main

import (
	"fmt"
	"os"

	"github.com/insighted4/correios-cep/pkg/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix        = "ADMIN"
	longDescription  = app.Description
	shortDescription = "admin"
)

var environment = os.Getenv("ENVIRONMENT")

func commandRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   shortDescription,
		Short: shortDescription,
		Long:  longDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(2)
		},
	}

	viper.AutomaticEnv()

	rootCmd.AddCommand(commandServe())
	rootCmd.AddCommand(newVersion(longDescription))

	return rootCmd
}

func main() {
	if err := commandRoot().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
