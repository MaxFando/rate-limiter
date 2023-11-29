package bucket

import (
	"context"
)

type Service interface {
	TryGetPermissionInLoginBucket(ctx context.Context, key string, limit int) bool
	TryGetPermissionInPasswordBucket(ctx context.Context, key string, limit int) bool
	ResetLoginBucket(ctx context.Context, login string) bool
	ResetIPBucket(ctx context.Context, ip string) bool
}

type UseCase struct {
	bucketService Service
}

func NewUseCase(bucketService Service) *UseCase {
	return &UseCase{bucketService: bucketService}
}

func (a *UseCase) Reset(ctx context.Context, login, ip string) (bool, bool, error) {
	isLoginReset := a.bucketService.ResetLoginBucket(ctx, login)
	isIPReset := a.bucketService.ResetIPBucket(ctx, ip)

	return isLoginReset, isIPReset, nil
}
