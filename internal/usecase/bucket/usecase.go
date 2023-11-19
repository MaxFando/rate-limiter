package bucket

import (
	"context"
	"errors"
	"github.com/MaxFando/rate-limiter/pkg/tracing"
)

type Service interface {
	TryGetPermissionInLoginBucket(ctx context.Context, key string, limit int) bool
	TryGetPermissionInPasswordBucket(ctx context.Context, key string, limit int) bool
	ResetLoginBucket(ctx context.Context, login string) bool
	ResetIpBucket(ctx context.Context, ip string) bool
}

type UseCase struct {
	bucketService Service
}

func NewUseCase(bucketService Service) *UseCase {
	return &UseCase{bucketService: bucketService}
}

func (a *UseCase) Reset(ctx context.Context, login, ip string) (bool, error) {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	isLoginReset := a.bucketService.ResetLoginBucket(ctx, login)
	if !isLoginReset {
		return false, errors.New("login bucket not reset")
	}

	isIpReset := a.bucketService.ResetIpBucket(ctx, ip)
	if !isIpReset {
		return false, errors.New("ip bucket not reset")
	}

	return true, nil
}
