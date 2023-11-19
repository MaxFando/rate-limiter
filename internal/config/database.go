package config

import "github.com/spf13/viper"

type Database struct {
	Host     string
	DbName   string
	Port     string
	User     string
	Password string
	SslMode  string
}

func newDatabase() Database {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_NAME", "postgres")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("SSL_MODE", "disable")

	return Database{
		Host:     viper.GetString("DB_HOST"),
		DbName:   viper.GetString("DB_NAME"),
		Port:     viper.GetString("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		SslMode:  viper.GetString("SSL_MODE"),
	}
}
