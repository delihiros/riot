package config

import (
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var config *viper.Viper

func GetConfig() *viper.Viper {
	return config
}

func ReadConfig(cfgFile string) {
	config = viper.New()
	if cfgFile != "" {
		config.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		config.AddConfigPath(home)
		config.SetConfigName(".riot")
	}
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
