package config

import "github.com/spf13/viper"

type appConfig struct {
	LogLevel string
}

func newAppConfig() appConfig {
	viper.SetDefault("LOGLEVEL", "debug")

	return appConfig{
		LogLevel: viper.GetString("LOGLEVEL"),
	}
}
