package grpcapi

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpcapi/blacklistpb"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/usecase/blacklist"
	"github.com/MaxFando/rate-limiter/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *BlacklistServer) GetIpList(ctx *blacklistpb.GetIpListRequest, stream blacklistpb.BlackListService_GetIpListServer) error {
	list, err := s.uc.GetIPList(context.Background())
	if err != nil {
		utils.Logger.Error("GetIpList error:", zap.Error(err))
		return err
	}

	var ipList []*blacklistpb.IpNetwork
	for _, _network := range list {
		ipList = append(ipList, &blacklistpb.IpNetwork{
			Ip:   _network.Ip.String(),
			Mask: _network.Mask.String(),
		})
	}

	response := &blacklistpb.GetIpListResponse{
		IpNetwork: ipList,
	}

	errStream := stream.Send(response)
	if errStream != nil {
		utils.Logger.Error("GetIpList error:", zap.Error(errStream))
		return status.Errorf(codes.Internal, "unexpected error: %v", errStream)
	}

	return nil
}
