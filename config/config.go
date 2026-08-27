package config

import (
	"os"

	"github.com/spf13/viper"
)

var config *viper.Viper

func GetConfig() *viper.Viper {
	return config
}

func ReadConfig(cfgFile string) error {
	config = viper.New()
	if err := config.BindEnv("riot.api_key", "RIOT_API_KEY"); err != nil {
		return err
	}
	if cfgFile != "" {
		config.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		config.AddConfigPath(home)
		config.SetConfigName(".riot")
	}
	err := config.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); cfgFile == "" && ok {
		return nil
	}
	return err
}
