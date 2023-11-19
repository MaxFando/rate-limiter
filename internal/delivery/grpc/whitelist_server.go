package grpc

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpc/whitelistpb"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/usecase/whitelist"
)

type WhitelistServer struct {
	whitelistpb.WhiteListServiceServer
	uc *whitelist.UseCase
}

func NewWhitelistServer(uc *whitelist.UseCase) *WhitelistServer {
	return &WhitelistServer{uc: uc}
}

func (s *WhitelistServer) AddIp(ctx context.Context, req *whitelistpb.AddIpRequest) (*whitelistpb.AddIpResponse, error) {
	request, err := network.NewIpNetwork(
		req.IpNetwork.Ip,
		req.IpNetwork.Mask,
	)

	err = s.uc.AddIP(ctx, request)
	if err != nil {
		return nil, err
	}

	return &whitelistpb.AddIpResponse{Ok: true}, nil
}

func (s *WhitelistServer) RemoveIp(ctx context.Context, req *whitelistpb.RemoveIPRequest) (*whitelistpb.RemoveIPResponse, error) {
	request, err := network.NewIpNetwork(req.IpNetwork.Ip, req.IpNetwork.Mask)
	if err != nil {
		return nil, err
	}

	err = s.uc.RemoveIP(ctx, request)
	if err != nil {
		return nil, err
	}

	return &whitelistpb.RemoveIPResponse{Ok: true}, nil
}

func (s *WhitelistServer) GetIpList(ctx context.Context, req *whitelistpb.GetIpListRequest) (*whitelistpb.GetIpListResponse, error) {
	list, err := s.uc.GetIPList(ctx)
	if err != nil {
		return nil, err
	}

	var ipList []*whitelistpb.IpNetwork
	for _, network := range list {
		ipList = append(ipList, &whitelistpb.IpNetwork{
			Ip:   network.Ip.String(),
			Mask: network.Mask.String(),
		})
	}

	return &whitelistpb.GetIpListResponse{IpNetwork: ipList}, nil
}
