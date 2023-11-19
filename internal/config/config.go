package config

import (
	"sync"

	"github.com/spf13/viper"
)

type config struct {
	IsDebug       bool
	IsDevelopment bool
	Listen        listen
	Server        server
	AppConfig     appConfig
	Database      Database
	Bucket        bucket
}

var (
	Config config
	once   sync.Once
)

func InitializeConfig() {
	once.Do(func() {
		Config = config{
			IsDebug:       viper.GetBool("IS_DEBUG"),
			IsDevelopment: viper.GetBool("IS_DEVELOP"),
			Listen:        newListen(),
			Server:        newServer(),
			AppConfig:     newAppConfig(),
			Database:      newDatabase(),
			Bucket:        newBucket(),
		}
	})
}
