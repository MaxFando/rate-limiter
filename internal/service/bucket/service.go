package bucket

import (
	"context"

	"github.com/MaxFando/rate-limiter/pkg/tracing"
)

type Repository interface {
	DeleteUnusedBucket(ctx context.Context)
	TryGetPermissionInBucket(ctx context.Context, requestValue string, limit int) bool
	ResetBucket(ctx context.Context, requestValue string) bool
}

type Service struct {
	ipBucketRepo       Repository
	loginBucketRepo    Repository
	passwordBucketRepo Repository
}

func NewService(
	ipBucketRepo Repository,
	loginBucketRepo Repository,
	passwordBucketRepo Repository,
) *Service {
	s := &Service{ipBucketRepo: ipBucketRepo, loginBucketRepo: loginBucketRepo, passwordBucketRepo: passwordBucketRepo}
	s.deleteUnusedBucket(context.Background())

	return s
}

func (s *Service) TryGetPermissionInLoginBucket(ctx context.Context, key string, limit int) bool {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return s.loginBucketRepo.TryGetPermissionInBucket(ctx, key, limit)
}

func (s *Service) TryGetPermissionInPasswordBucket(ctx context.Context, key string, limit int) bool {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return s.passwordBucketRepo.TryGetPermissionInBucket(ctx, key, limit)
}

func (s *Service) ResetLoginBucket(ctx context.Context, login string) bool {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return s.loginBucketRepo.ResetBucket(ctx, login)
}

func (s *Service) ResetIpBucket(ctx context.Context, ip string) bool {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return s.ipBucketRepo.ResetBucket(ctx, ip)
}

func (s *Service) deleteUnusedBucket(ctx context.Context) {
	go s.ipBucketRepo.DeleteUnusedBucket(ctx)
	go s.loginBucketRepo.DeleteUnusedBucket(ctx)
	go s.passwordBucketRepo.DeleteUnusedBucket(ctx)
}
