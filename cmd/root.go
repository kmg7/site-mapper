package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: version,
	Use:     "site-mapper",
	Short:   "A basic website map builder",
	Long:    `Web site map builder to start building site maps provide an url with -u flag`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
