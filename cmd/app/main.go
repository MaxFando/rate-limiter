package main

import (
	"context"
	"fmt"
	"github.com/MaxFando/rate-limiter/internal/app/http"
	"github.com/MaxFando/rate-limiter/internal/config"
	"github.com/MaxFando/rate-limiter/internal/providers"
	"github.com/MaxFando/rate-limiter/internal/store/postgres"
	"github.com/MaxFando/rate-limiter/pkg/tracing"
	"github.com/MaxFando/rate-limiter/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	utils.InitializeLogger()

	ctx, cancel := context.WithCancel(context.Background())
	config.InitializeConfig()

	closer, err := tracing.New("discounts")
	if err != nil {
		utils.Logger.Error("Не удалось подключиться к jaeger", zap.Error(err))
	}
	defer closer.Close()

	postgresDB := postgres.New()
	err = postgresDB.Connect(ctx, "default", config.Config.Database)
	if err != nil {
		panic(err)
	}

	repositoryProvider := providers.NewRepositoryProvider(postgresDB)
	repositoryProvider.BootPrefixProviderContract()

	serviceProvider := providers.NewServiceProvider()
	serviceProvider.RegisterDependencies(repositoryProvider)

	ctx = context.WithValue(ctx, providers.ServiceProviderKey, serviceProvider)

	httpServer := http.NewHttpServer(http.NewHandler(ctx), ":"+config.Config.Listen.Port)
	httpServer.Serve()
	utils.Logger.Info("Приложение стартовало в режиме", zap.String("log_level", config.Config.AppConfig.LogLevel))
	utils.Logger.Info("На порту " + config.Config.Listen.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		utils.Logger.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		utils.Logger.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	cancel()
	utils.Logger.Info("Приложение завершило работу")
}
