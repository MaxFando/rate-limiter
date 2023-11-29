package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/MaxFando/rate-limiter/internal/config"
	"github.com/MaxFando/rate-limiter/internal/store/postgres/plugins/logger"
	"github.com/MaxFando/rate-limiter/pkg/utils"
)

type Client struct {
	Connection map[string]*gorm.DB
}

func New() *Client {
	return &Client{
		Connection: map[string]*gorm.DB{},
	}
}

func (c *Client) Connect(ctx context.Context, name string, cfg config.Database) error {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SslMode)

	connection, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.NewGormZapLoggerWrapper(utils.Logger.GetLogger().Sugar()).LogMode(gormLogger.Info),
	})

	if err != nil {
		utils.Logger.Fatal(err)
		return err
	}

	sqlDB, err := connection.DB()
	if err != nil {
		utils.Logger.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	c.Connection[name] = connection

	return nil
}
