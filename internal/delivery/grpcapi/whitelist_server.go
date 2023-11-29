package grpcapi

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpcapi/whitelistpb"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/usecase/whitelist"
	"github.com/MaxFando/rate-limiter/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WhitelistServer struct {
	whitelistpb.WhiteListServiceServer
	uc *whitelist.UseCase
}

func NewWhiteListServer(uc *whitelist.UseCase) *WhitelistServer {
	return &WhitelistServer{uc: uc}
}

func (s *WhitelistServer) AddIP(ctx context.Context, req *whitelistpb.AddIPRequest) (*whitelistpb.AddIPResponse, error) {
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

	return &whitelistpb.AddIPResponse{Ok: true}, nil
}

func (s *WhitelistServer) RemoveIP(ctx context.Context, req *whitelistpb.RemoveIPRequest) (*whitelistpb.RemoveIPResponse, error) {
	request, err := network.NewIPNetwork(req.IpNetwork.Ip, req.IpNetwork.Mask)
	if err != nil {
		return nil, err
	}

	err = s.uc.RemoveIP(ctx, request)
	if err != nil {
		return nil, err
	}

	return &whitelistpb.RemoveIPResponse{Ok: true}, nil
}

func (s *WhitelistServer) GetIPList(req *whitelistpb.GetIpListRequest, stream whitelistpb.WhiteListService_GetIpListServer) error {
	list, err := s.uc.GetIPList(context.Background())
	if err != nil {
		utils.Logger.Error("GetIPList error:", zap.Error(err))
		return err
	}

	ipList := make([]*whitelistpb.IPNetwork, 0, len(list))
	for _, _network := range list {
		ipList = append(ipList, &whitelistpb.IPNetwork{
			Ip:   _network.IP.String(),
			Mask: _network.Mask.String(),
		})
	}

	response := &whitelistpb.GetIPListResponse{
		IpNetwork: ipList,
	}

	errStream := stream.Send(response)
	if errStream != nil {
		utils.Logger.Error("GetIPList error:", zap.Error(errStream))
		return status.Errorf(codes.Internal, "unexpected error: %v", errStream)
	}

	return nil
}
