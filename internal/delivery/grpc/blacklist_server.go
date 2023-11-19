package grpc

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpc/blacklistpb"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/usecase/blacklist"
)

type BlacklistServer struct {
	blacklistpb.BlackListServiceServer
	uc *blacklist.UseCase
}

func NewBlacklistServer(uc *blacklist.UseCase) *BlacklistServer {
	return &BlacklistServer{uc: uc}
}

func (s *BlacklistServer) AddIp(ctx context.Context, req *blacklistpb.AddIpRequest) (*blacklistpb.AddIpResponse, error) {
	request, err := network.NewIpNetwork(
		req.IpNetwork.Ip,
		req.IpNetwork.Mask,
	)

	err = s.uc.AddIP(ctx, request)
	if err != nil {
		return nil, err
	}

	return &blacklistpb.AddIpResponse{Ok: true}, nil
}

func (s *BlacklistServer) RemoveIp(ctx context.Context, req *blacklistpb.RemoveIPRequest) (*blacklistpb.RemoveIPResponse, error) {
	request, err := network.NewIpNetwork(req.IpNetwork.Ip, req.IpNetwork.Mask)
	if err != nil {
		return nil, err
	}

	err = s.uc.RemoveIP(ctx, request)
	if err != nil {
		return nil, err
	}

	return &blacklistpb.RemoveIPResponse{Ok: true}, nil
}

func (s *BlacklistServer) GetIpList(ctx context.Context, req *blacklistpb.GetIpListRequest) (*blacklistpb.GetIpListResponse, error) {
	list, err := s.uc.GetIPList(ctx)
	if err != nil {
		return nil, err
	}

	var ipList []*blacklistpb.IpNetwork
	for _, network := range list {
		ipList = append(ipList, &blacklistpb.IpNetwork{
			Ip:   network.Ip.String(),
			Mask: network.Mask.String(),
		})
	}

	return &blacklistpb.GetIpListResponse{IpNetwork: ipList}, nil
}
