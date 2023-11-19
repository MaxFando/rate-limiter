package http

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/auth"
	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/blacklist"
	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/bucket"
	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/whitelist"
	"github.com/MaxFando/rate-limiter/internal/providers"
	authUC "github.com/MaxFando/rate-limiter/internal/usecase/auth"
	blacklistUC "github.com/MaxFando/rate-limiter/internal/usecase/blacklist"
	bucketUC "github.com/MaxFando/rate-limiter/internal/usecase/bucket"
	whiteListUC "github.com/MaxFando/rate-limiter/internal/usecase/whitelist"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func NewRouter(ctx context.Context, handler *echo.Echo) *echo.Echo {
	handler.GET("/metrics", echoprometheus.NewHandler())

	apiGroup := handler.Group("/api")
	v1Group := apiGroup.Group("/v1")

	serviceProvider := ctx.Value(providers.ServiceProviderKey).(*providers.ServiceProvider)

	authUseCase := authUC.NewUseCase(serviceProvider.BlacklistService, serviceProvider.WhitelistService, serviceProvider.BucketService)
	blackListUseCase := blacklistUC.NewUseCase(serviceProvider.BlacklistService)
	whiteListUseCase := whiteListUC.NewUseCase(serviceProvider.WhitelistService)
	bucketUseCase := bucketUC.NewUseCase(serviceProvider.BucketService)

	authControllerV1 := auth.NewAuthController(authUseCase)
	blackListControllerV1 := blacklist.NewController(blackListUseCase)
	whiteListControllerV1 := whitelist.NewController(whiteListUseCase)
	bucketControllerV1 := bucket.NewController(bucketUseCase)

	v1Group.POST("/auth", authControllerV1.TryAuthorization)
	v1Group.DELETE("/auth/bucket", bucketControllerV1.ResetBucket)
	v1Group.POST("/auth/blacklist", blackListControllerV1.AddIP)
	v1Group.DELETE("/auth/blacklist", blackListControllerV1.RemoveIP)
	v1Group.GET("/auth/blacklist", blackListControllerV1.GetIPList)
	v1Group.POST("/auth/whitelist", whiteListControllerV1.AddIP)
	v1Group.DELETE("/auth/whitelist", whiteListControllerV1.RemoveIP)
	v1Group.GET("/auth/whitelist", whiteListControllerV1.GetIPList)

	return handler
}
