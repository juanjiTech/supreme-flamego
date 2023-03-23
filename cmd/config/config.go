package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"supreme-flamego/conf"
)

var (
	configYml string
	forceGen  bool
	StartCmd  = &cobra.Command{
		Use:     "config",
		Short:   "Generate config file",
		Example: "mod config -p config/config.yaml -f",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Generate config...")
			err := load()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "path", "p", "config/config.yaml", "Generate config in provided path")
	StartCmd.PersistentFlags().BoolVarP(&forceGen, "force", "f", false, "Force generate config in provided path")
}

func load() error {
	return conf.GenConfig(configYml, forceGen)
}
