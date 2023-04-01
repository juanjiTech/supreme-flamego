package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"supreme-flamego/pkg/colorful"
	"supreme-flamego/pkg/fsx"
)

var SysVersion = "dev"

var serveConfig *GlobalConfig

func LoadConfig(configPath ...string) {
	if len(configPath) == 0 || configPath[0] == "" {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	} else {
		viper.SetConfigFile(configPath[0])
	}

	serveConfig = new(GlobalConfig)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config Read failed: " + err.Error())
		os.Exit(1)
	}
	err = viper.Unmarshal(serveConfig)
	if err != nil {
		fmt.Println("Config Unmarshal failed: " + err.Error())
		os.Exit(1)
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config fileHandle changed: ", e.Name)
		_ = viper.ReadInConfig()
		err = viper.Unmarshal(serveConfig)
		if err != nil {
			fmt.Println("New Config fileHandle Parse Failed: ", e.Name)
			return
		}
	})
	viper.WatchConfig()
}

func GenYamlConfig(path string, force bool) error {
	if !fsx.FileExist(path) || force {
		data, _ := yaml.Marshal(&GlobalConfig{MODE: "debug"})
		err := os.WriteFile(path, data, 0644)
		if err != nil {
			return errors.New(colorful.Red("Generate file with error: " + err.Error()))
		}
		fmt.Println(colorful.Green("Config file `config.yaml` generate success in " + path))
	} else {
		return errors.New(colorful.Red(path + " already exist, use -f to Force coverage"))
	}
	return nil
}

func GetConfig() *GlobalConfig {
	return serveConfig
}
