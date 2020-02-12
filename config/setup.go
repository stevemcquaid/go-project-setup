package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func SetupConfig(input interface{}) error {
	// Add config paths to search
	viper.AddConfigPath("/etc/hello")
	viper.AddConfigPath("$HOME/.hello")
	viper.AddConfigPath("/")
	viper.AddConfigPath(".")
	// Add in config file
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file, got error: %s", err)
	}
	err := viper.Unmarshal(&input)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	return nil
}
