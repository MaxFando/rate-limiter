package config

import "github.com/spf13/viper"

type server struct {
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	ServerType   string
}

func newServer() server {
	viper.SetDefault("READ_TIMEOUT", 10)
	viper.SetDefault("WRITE_TIMEOUT", 10)
	viper.SetDefault("IDLE_TIMEOUT", 100)
	viper.SetDefault("SERVER_TYPE", "http")

	return server{
		ReadTimeout:  viper.GetInt("READ_TIMEOUT"),
		WriteTimeout: viper.GetInt("WRITE_TIMEOUT"),
		IdleTimeout:  viper.GetInt("IDLE_TIMEOUT"),
		ServerType:   viper.GetString("SERVER_TYPE"),
	}
}
