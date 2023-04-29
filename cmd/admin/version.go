package main

import (
	"fmt"
	"runtime"

	"github.com/insighted4/correios-cep/pkg/version"
	"github.com/spf13/cobra"
)

func newVersion(description string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version and exit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(description)
			fmt.Printf("Go Version: %s\n", runtime.Version())
			fmt.Printf("Go OS/ARCH: %s %s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Printf("Build Time: %s\n", version.BuildTime)
			fmt.Printf("Commit: %s\n", version.CommitHash)
			fmt.Printf("Version: %s\n", version.Version)
		},
	}
}
