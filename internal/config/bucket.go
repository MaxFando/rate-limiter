package config

import "github.com/spf13/viper"

type bucket struct {
	IPLimit             int
	LoginLimit          int
	PasswordLimit       int
	ResetBucketInterval int
}

func newBucket() bucket {
	viper.SetDefault("IP_LIMIT", 1000)
	viper.SetDefault("LOGIN_LIMIT", 10)
	viper.SetDefault("PASSWORD_LIMIT", 100)
	viper.SetDefault("RESET_BUCKET_INTERVAL", 60)

	return bucket{
		IPLimit:             viper.GetInt("IP_LIMIT"),
		LoginLimit:          viper.GetInt("LOGIN_LIMIT"),
		PasswordLimit:       viper.GetInt("PASSWORD_LIMIT"),
		ResetBucketInterval: viper.GetInt("RESET_BUCKET_INTERVAL"),
	}
}
