package grpcapi

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/delivery/grpcapi/bucketpb"
	"github.com/MaxFando/rate-limiter/internal/usecase/bucket"
)

type BucketServer struct {
	bucketpb.BucketServiceServer
	uc *bucket.UseCase
}

func NewBucketServer(uc *bucket.UseCase) *BucketServer {
	return &BucketServer{uc: uc}
}

func (s *BucketServer) ResetBucket(ctx context.Context, req *bucketpb.ResetBucketRequest) (*bucketpb.ResetBucketResponse, error) {
	isLoginReset, isIPReset, err := s.uc.Reset(ctx, req.Request.Login, req.Request.Ip)
	if err != nil {
		return nil, err
	}

	return &bucketpb.ResetBucketResponse{ResetIp: isIPReset, ResetLogin: isLoginReset}, nil
}
