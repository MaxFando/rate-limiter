package grpcapi

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/config"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpcapi"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpcapi/authpb"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpcapi/blacklistpb"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpcapi/bucketpb"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpcapi/whitelistpb"
	"github.com/MaxFando/rate-limiter/internal/providers"
	"github.com/MaxFando/rate-limiter/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

type Server struct {
	blackListServer     *grpcapi.BlacklistServer
	whiteListServer     *grpcapi.WhitelistServer
	bucketServer        *grpcapi.BucketServer
	authorizationServer *grpcapi.AuthServer
	grpcServer          *grpc.Server

	errors chan error
}

func NewServer(ctx context.Context) *Server {
	grpcServer := grpc.NewServer()

	useCaseProvider := ctx.Value(providers.UseCaseProviderKey).(*providers.UseCaseProvider)
	blackListServer := grpcapi.NewBlacklistServer(useCaseProvider.BlackListUseCase)
	whiteListServer := grpcapi.NewWhiteListServer(useCaseProvider.WhiteListUseCase)
	bucketServer := grpcapi.NewBucketServer(useCaseProvider.BucketUseCase)
	authorizationServer := grpcapi.NewAuthServer(useCaseProvider.AuthUseCase)

	return &Server{
		blackListServer:     blackListServer,
		whiteListServer:     whiteListServer,
		bucketServer:        bucketServer,
		authorizationServer: authorizationServer,
		grpcServer:          grpcServer,

		errors: make(chan error),
	}
}

func (s *Server) Serve() {
	utils.Logger.Info("Start GRPC Server")

	listener, _ := net.Listen("tcp", config.Config.Listen.BindIP+":"+config.Config.Listen.Port)

	blacklistpb.RegisterBlackListServiceServer(s.grpcServer, s.blackListServer)
	whitelistpb.RegisterWhiteListServiceServer(s.grpcServer, s.whiteListServer)
	bucketpb.RegisterBucketServiceServer(s.grpcServer, s.bucketServer)
	authpb.RegisterAuthorizationServer(s.grpcServer, s.authorizationServer)

	reflection.Register(s.grpcServer)
	go func() {
		s.errors <- s.grpcServer.Serve(listener)
		close(s.errors)
	}()
}

func (s *Server) Shutdown(c chan os.Signal) {
	utils.Logger.Info("Service is stop")
	s.grpcServer.GracefulStop()
}

func (s *Server) Notify() <-chan error {
	return s.errors
}
