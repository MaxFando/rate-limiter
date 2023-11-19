package blacklist

import (
	"context"

	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/pkg/tracing"
)

type Service interface {
	AddIP(ctx context.Context, network network.IpNetwork) error
	RemoveIP(ctx context.Context, network network.IpNetwork) error
	GetIPList(ctx context.Context) ([]network.IpNetwork, error)
}

type UseCase struct {
	blackListService Service
}

func NewUseCase(blackListService Service) *UseCase {
	return &UseCase{blackListService: blackListService}
}

func (u *UseCase) AddIP(ctx context.Context, network network.IpNetwork) error {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return u.blackListService.AddIP(ctx, network)
}

func (u *UseCase) RemoveIP(ctx context.Context, network network.IpNetwork) error {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return u.blackListService.RemoveIP(ctx, network)
}

func (u *UseCase) GetIPList(ctx context.Context) ([]network.IpNetwork, error) {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return u.blackListService.GetIPList(ctx)
}
