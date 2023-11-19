package grpcapi

import (
	"github.com/MaxFando/rate-limiter/internal/config"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpc/authpb"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpc/blacklistpb"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpc/bucketpb"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpc/whitelistpb"
	"github.com/MaxFando/rate-limiter/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

type Server struct {
	blackListServer     blacklistpb.BlackListServiceServer
	whiteListServer     whitelistpb.WhiteListServiceServer
	bucketServer        bucketpb.BucketServiceServer
	authorizationServer authpb.AuthorizationServer
	grpcServer          *grpc.Server

	errors chan error
}

func NewServer(
	blackListServer blacklistpb.BlackListServiceServer,
	whiteListServer whitelistpb.WhiteListServiceServer,
	bucketServer bucketpb.BucketServiceServer,
	authorizationServer authpb.AuthorizationServer,
) *Server {
	grpcServer := grpc.NewServer()

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

	// todo: не игнорировать ошибку
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
