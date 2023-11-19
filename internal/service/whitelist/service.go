package whitelist

import (
	"context"

	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/pkg/tracing"
	"github.com/MaxFando/rate-limiter/pkg/utils"
)

type Store interface {
	Add(ctx context.Context, prefix string, mask string) error
	Remove(ctx context.Context, prefix string, mask string) error
	Get(ctx context.Context) ([]network.IpNetwork, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) AddIP(ctx context.Context, network network.IpNetwork) error {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	prefix, err := utils.GetPrefix(network.Ip.String(), network.Mask.String())
	if err != nil {
		return err
	}
	err = s.store.Add(ctx, prefix, network.Mask.String())
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveIP(ctx context.Context, network network.IpNetwork) error {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	prefix, err := utils.GetPrefix(network.Ip.String(), network.Mask.String())
	if err != nil {
		return err
	}
	err = s.store.Remove(ctx, prefix, network.Mask.String())
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetIPList(ctx context.Context) ([]network.IpNetwork, error) {
	span, ctx := tracing.CreateChildSpanWithFuncName(ctx)
	defer span.Finish()

	ipList, err := s.store.Get(ctx)
	if err != nil {
		return nil, err
	}
	return ipList, nil
}
