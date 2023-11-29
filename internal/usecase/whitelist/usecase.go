package whitelist

import (
	"context"

	"github.com/MaxFando/rate-limiter/internal/domain/network"
)

type Service interface {
	AddIP(ctx context.Context, network network.IPNetwork) error
	RemoveIP(ctx context.Context, network network.IPNetwork) error
	GetIPList(ctx context.Context) ([]network.IPNetwork, error)
}

type UseCase struct {
	whiteListService Service
}

func NewUseCase(whiteListService Service) *UseCase {
	return &UseCase{whiteListService: whiteListService}
}

func (u *UseCase) AddIP(ctx context.Context, network network.IPNetwork) error {
	return u.whiteListService.AddIP(ctx, network)
}

func (u *UseCase) RemoveIP(ctx context.Context, network network.IPNetwork) error {
	return u.whiteListService.RemoveIP(ctx, network)
}

func (u *UseCase) GetIPList(ctx context.Context) ([]network.IPNetwork, error) {
	return u.whiteListService.GetIPList(ctx)
}
