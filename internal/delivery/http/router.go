package http

import (
	"context"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"

	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/auth"
	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/blacklist"
	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/bucket"
	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/whitelist"
	"github.com/MaxFando/rate-limiter/internal/providers"
)

func NewRouter(ctx context.Context, handler *echo.Echo) *echo.Echo {
	handler.GET("/metrics", echoprometheus.NewHandler())

	apiGroup := handler.Group("/api")
	v1Group := apiGroup.Group("/v1")

	useCaseProvider := ctx.Value(providers.UseCaseProviderKey).(*providers.UseCaseProvider)

	authControllerV1 := auth.NewAuthController(useCaseProvider.AuthUseCase)
	blackListControllerV1 := blacklist.NewController(useCaseProvider.BlackListUseCase)
	whiteListControllerV1 := whitelist.NewController(useCaseProvider.WhiteListUseCase)
	bucketControllerV1 := bucket.NewController(useCaseProvider.BucketUseCase)

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
