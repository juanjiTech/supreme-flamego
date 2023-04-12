package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"supreme-flamego/cmd/config"
	"supreme-flamego/cmd/create"
	"supreme-flamego/cmd/server"
)

var rootCmd = &cobra.Command{
	Use:          "mod",
	Short:        "mod",
	SilenceUsage: true,
	Long:         `mod`,
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(create.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
