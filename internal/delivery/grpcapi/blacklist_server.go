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

func (s *BlacklistServer) AddIP(ctx context.Context, req *blacklistpb.AddIPRequest) (*blacklistpb.AddIPResponse, error) {
	request, err := network.NewIPNetwork(
		req.IpNetwork.Ip,
		req.IpNetwork.Mask,
	)
	if err != nil {
		return nil, err
	}

	err = s.uc.AddIP(ctx, request)
	if err != nil {
		return nil, err
	}

	return &blacklistpb.AddIPResponse{Ok: true}, nil
}

func (s *BlacklistServer) RemoveIP(ctx context.Context, req *blacklistpb.RemoveIPRequest) (*blacklistpb.RemoveIPResponse, error) {
	request, err := network.NewIPNetwork(req.IpNetwork.Ip, req.IpNetwork.Mask)
	if err != nil {
		return nil, err
	}

	err = s.uc.RemoveIP(ctx, request)
	if err != nil {
		return nil, err
	}

	return &blacklistpb.RemoveIPResponse{Ok: true}, nil
}

func (s *BlacklistServer) GetIPList(ctx *blacklistpb.GetIPListRequest, stream blacklistpb.BlackListService_GetIpListServer) error {
	list, err := s.uc.GetIPList(context.Background())
	if err != nil {
		utils.Logger.Error("GetIPList error:", zap.Error(err))
		return err
	}

	ipList := make([]*blacklistpb.IPNetwork, 0, len(list))
	for _, _network := range list {
		ipList = append(ipList, &blacklistpb.IPNetwork{
			Ip:   _network.IP.String(),
			Mask: _network.Mask.String(),
		})
	}

	response := &blacklistpb.GetIPListResponse{
		IpNetwork: ipList,
	}

	errStream := stream.Send(response)
	if errStream != nil {
		utils.Logger.Error("GetIPList error:", zap.Error(errStream))
		return status.Errorf(codes.Internal, "unexpected error: %v", errStream)
	}

	return nil
}
