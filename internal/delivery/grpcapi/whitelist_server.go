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

func (s *WhitelistServer) GetIpList(req *whitelistpb.GetIpListRequest, stream whitelistpb.WhiteListService_GetIpListServer) error {
	list, err := s.uc.GetIPList(context.Background())
	if err != nil {
		utils.Logger.Error("GetIpList error:", zap.Error(err))
		return err
	}

	var ipList []*whitelistpb.IpNetwork
	for _, _network := range list {
		ipList = append(ipList, &whitelistpb.IpNetwork{
			Ip:   _network.Ip.String(),
			Mask: _network.Mask.String(),
		})
	}

	response := &whitelistpb.GetIpListResponse{
		IpNetwork: ipList,
	}

	errStream := stream.Send(response)
	if errStream != nil {
		utils.Logger.Error("GetIpList error:", zap.Error(errStream))
		return status.Errorf(codes.Internal, "unexpected error: %v", errStream)
	}

	return nil
}
