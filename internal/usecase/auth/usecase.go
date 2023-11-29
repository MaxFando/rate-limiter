package auth

import (
	"context"

	"github.com/MaxFando/rate-limiter/internal/config"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/service/blacklist"
	"github.com/MaxFando/rate-limiter/internal/service/bucket"
	"github.com/MaxFando/rate-limiter/internal/service/whitelist"
	"github.com/MaxFando/rate-limiter/pkg/utils"
)

type UseCase struct {
	blackListService *blacklist.Service
	whiteListService *whitelist.Service
	bucketService    *bucket.Service
}

func NewUseCase(blackListService *blacklist.Service, whiteListService *whitelist.Service, bucketService *bucket.Service) *UseCase {
	return &UseCase{blackListService: blackListService, whiteListService: whiteListService, bucketService: bucketService}
}

func (uc *UseCase) TryAuthorization(ctx context.Context, request network.Request) (bool, error) {
	utils.Logger.Info("Check ip in blacklist")
	ipNetworkList, err := uc.blackListService.GetIPList(ctx)
	if err != nil {
		return false, err
	}
	isIPInBlackList, err := uc.checkIPByNetworkList(ctx, request.IP.String(), ipNetworkList)
	if err != nil {
		return false, err
	}
	if isIPInBlackList {
		return false, nil
	}

	utils.Logger.Info("Check ip in whitelist")
	ipNetworkList, err = uc.whiteListService.GetIPList(ctx)
	if err != nil {
		return false, err
	}
	isIPInWhiteList, err := uc.checkIPByNetworkList(ctx, request.IP.String(), ipNetworkList)
	if err != nil {
		return false, err
	}
	if isIPInWhiteList {
		return true, nil
	}

	utils.Logger.Info("Check ip in bucketService")
	isAllow := true
	allow := uc.bucketService.TryGetPermissionInLoginBucket(ctx, request.IP.String(), config.Config.Bucket.IPLimit)
	if !allow {
		isAllow = allow
	}

	utils.Logger.Info("Check password in bucketService")
	allow = uc.bucketService.TryGetPermissionInPasswordBucket(ctx, request.Password, config.Config.Bucket.PasswordLimit)
	if !allow {
		isAllow = allow
	}

	utils.Logger.Info("Check login in bucketService")
	allow = uc.bucketService.TryGetPermissionInLoginBucket(ctx, request.Login, config.Config.Bucket.LoginLimit)
	if !allow {
		isAllow = allow
	}

	return isAllow, nil
}

func (uc *UseCase) checkIPByNetworkList(ctx context.Context, ip string, ipNetworkList []network.IPNetwork) (bool, error) {
	for i := range ipNetworkList {
		prefix, err := utils.GetPrefix(ip, ipNetworkList[i].Mask.String())
		if err != nil {
			return false, err
		}

		if prefix == ipNetworkList[i].IP.String() {
			return true, nil
		}
	}
	return false, nil
}
