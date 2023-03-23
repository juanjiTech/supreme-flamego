package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"supreme-flamego/pkg/colorful"
	"supreme-flamego/pkg/fs"
)

var SysVersion = "dev"

var serveConfig *GlobalConfig

func LoadConfig(configYml string) {
	if !fs.FileExist(configYml) {
		fmt.Println("cannot find config file")
		os.Exit(1)
	}
	serveConfig = new(GlobalConfig)
	viper.SetConfigFile(configYml)
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

func GenConfig(configYml string, force bool) error {
	if !fs.FileExist(configYml) || force {
		data, _ := yaml.Marshal(&GlobalConfig{MODE: "debug"})
		err := os.WriteFile(configYml, data, 0644)
		if err != nil {
			return errors.New(colorful.Red("Generate file with error: " + err.Error()))
		}
		fmt.Println(colorful.Green("config file `config.yaml` generate success in " + configYml))
	} else {
		return errors.New(colorful.Red(configYml + " already exist, use -f to Force coverage"))
	}
	return nil
}

func GetConfig() *GlobalConfig {
	return serveConfig
}
