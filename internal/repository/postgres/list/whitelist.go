package list

import (
	"context"
	"fmt"

	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/store/postgres"
)

const (
	isIPExistWl = `SELECT exists(SELECT 1 FROM whitelist WHERE prefix = $1 AND mask = $2)`
	insertIPWl  = `INSERT INTO whitelist (prefix, mask) VALUES ($1, $2)`
	deleteIPWl  = `DELETE FROM whitelist WHERE prefix = $1 AND mask = $2`
	getIPListWl = `SELECT prefix, mask from whitelist`
)

type WhiteListRepository struct {
	client *postgres.Client
}

func NewWhiteListRepository(client *postgres.Client) *WhiteListRepository {
	return &WhiteListRepository{client: client}
}

func (r *WhiteListRepository) Add(ctx context.Context, prefix string, mask string) error {
	var isExist bool

	err := r.client.Connection["default"].WithContext(ctx).Raw(isIPExistWl, prefix, mask).Scan(&isExist).Error
	if err != nil {
		return err
	}

	if isExist {
		return fmt.Errorf("this ip network already exist")
	}

	err = r.client.Connection["default"].WithContext(ctx).Exec(insertIPWl, prefix, mask).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *WhiteListRepository) Remove(ctx context.Context, prefix string, mask string) error {
	return r.client.Connection["default"].WithContext(ctx).Exec(deleteIPWl, prefix, mask).Error
}

func (r *WhiteListRepository) Get(ctx context.Context) ([]network.IPNetwork, error) {
	ipNetworkList := make([]network.IPNetwork, 0, 5)
	err := r.client.Connection["default"].WithContext(ctx).Raw(getIPListWl).Scan(&ipNetworkList).Error
	if err != nil {
		return nil, err
	}

	return ipNetworkList, nil
}
