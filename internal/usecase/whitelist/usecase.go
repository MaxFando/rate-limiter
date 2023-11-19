package whitelist

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
	whiteListService Service
}

func NewUseCase(whiteListService Service) *UseCase {
	return &UseCase{whiteListService: whiteListService}
}

func (u *UseCase) AddIP(ctx context.Context, network network.IpNetwork) error {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return u.whiteListService.AddIP(ctx, network)
}

func (u *UseCase) RemoveIP(ctx context.Context, network network.IpNetwork) error {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return u.whiteListService.RemoveIP(ctx, network)
}

func (u *UseCase) GetIPList(ctx context.Context) ([]network.IpNetwork, error) {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	return u.whiteListService.GetIPList(ctx)
}
