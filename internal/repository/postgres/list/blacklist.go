package list

import (
	"context"
	"fmt"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/store/postgres"
)

const (
	isIPExistBl = `SELECT exists(SELECT 1 FROM public.blacklist WHERE prefix = $1 AND mask = $2)`
	insertIPBl  = `INSERT INTO public.blacklist (prefix, mask) VALUES ($1, $2)`
	deleteIPBl  = `DELETE FROM public.blacklist WHERE prefix = $1 AND mask = $2`
	getIPListBl = `SELECT prefix, mask from public.blacklist`
)

type BlackListRepository struct {
	client *postgres.Client
}

func NewBlackListRepository(client *postgres.Client) *BlackListRepository {
	return &BlackListRepository{client: client}
}

func (r *BlackListRepository) Add(ctx context.Context, prefix string, mask string) error {
	var isExist bool

	err := r.client.Connection["default"].WithContext(ctx).Raw(isIPExistBl, prefix, mask).Scan(&isExist).Error
	if err != nil {
		return err
	}

	if isExist {
		return fmt.Errorf("this ip network already exist")
	}
	err = r.client.Connection["default"].WithContext(ctx).Exec(insertIPBl, prefix, mask).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *BlackListRepository) Remove(ctx context.Context, prefix string, mask string) error {
	return r.client.Connection["default"].WithContext(ctx).Exec(deleteIPBl, prefix, mask).Error
}

func (r *BlackListRepository) Get(ctx context.Context) ([]network.IpNetwork, error) {
	ipNetworkList := make([]network.IpNetwork, 0, 5)
	err := r.client.Connection["default"].WithContext(ctx).Raw(getIPListBl).Scan(&ipNetworkList).Error
	if err != nil {
		return nil, err
	}
	return ipNetworkList, nil
}
