package providers

import (
	authUC "github.com/MaxFando/rate-limiter/internal/usecase/auth"
	blacklistUC "github.com/MaxFando/rate-limiter/internal/usecase/blacklist"
	bucketUC "github.com/MaxFando/rate-limiter/internal/usecase/bucket"
	whiteListUC "github.com/MaxFando/rate-limiter/internal/usecase/whitelist"
)

type UseCaseProvider struct {
	AuthUseCase      *authUC.UseCase
	BlackListUseCase *blacklistUC.UseCase
	WhiteListUseCase *whiteListUC.UseCase
	BucketUseCase    *bucketUC.UseCase
}

func NewUseCaseProvider() *UseCaseProvider {
	return &UseCaseProvider{}
}

func (ucp *UseCaseProvider) RegisterDependencies(serviceProvider *ServiceProvider) {
	ucp.AuthUseCase = authUC.NewUseCase(serviceProvider.BlacklistService, serviceProvider.WhitelistService, serviceProvider.BucketService)
	ucp.BlackListUseCase = blacklistUC.NewUseCase(serviceProvider.BlacklistService)
	ucp.WhiteListUseCase = whiteListUC.NewUseCase(serviceProvider.WhitelistService)
	ucp.BucketUseCase = bucketUC.NewUseCase(serviceProvider.BucketService)
}
