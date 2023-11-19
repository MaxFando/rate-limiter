package config

import "github.com/spf13/viper"

type listen struct {
	Type   string
	BindIP string
	Port   string
}

func newListen() listen {
	viper.SetDefault("LISTEN_TYPE", "port")
	viper.SetDefault("BIND_IP", "0.0.0.0")
	viper.SetDefault("HTTP_PORT", "8080")

	return listen{
		Type:   viper.GetString("LISTEN_TYPE"),
		BindIP: viper.GetString("BIND_IP"),
		Port:   viper.GetString("HTTP_PORT"),
	}
}
