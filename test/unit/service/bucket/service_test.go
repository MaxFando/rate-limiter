package bucket

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/service/bucket"
	mocks "github.com/MaxFando/rate-limiter/mocks/service/bucket"
	"github.com/MaxFando/rate-limiter/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

var (
	mockIpBucketRepo       *mocks.Repository
	mockLoginBucketRepo    *mocks.Repository
	mockPasswordBucketRepo *mocks.Repository
)

func TestMain(m *testing.M) {
	utils.InitializeLogger()

	os.Exit(m.Run())
}

func TestService_TryGetPermissionInLoginBucket(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoginBucketRepo = new(mocks.Repository)
		mockIpBucketRepo = new(mocks.Repository)
		mockPasswordBucketRepo = new(mocks.Repository)

		key := "login"
		limit := 10
		mockLoginBucketRepo.On("TryGetPermissionInBucket", mock.Anything, key, limit).Return(true)
		s := bucket.NewService(mockIpBucketRepo, mockLoginBucketRepo, mockPasswordBucketRepo)

		allow := s.TryGetPermissionInLoginBucket(context.TODO(), key, limit)
		assert.True(t, allow)
	})
}

func TestService_TryGetPermissionInPasswordBucket(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoginBucketRepo = new(mocks.Repository)
		mockIpBucketRepo = new(mocks.Repository)
		mockPasswordBucketRepo = new(mocks.Repository)

		key := "password"
		limit := 10
		mockPasswordBucketRepo.On("TryGetPermissionInBucket", mock.Anything, key, limit).Return(true)
		s := bucket.NewService(mockIpBucketRepo, mockPasswordBucketRepo, mockPasswordBucketRepo)

		allow := s.TryGetPermissionInPasswordBucket(context.TODO(), key, limit)
		assert.True(t, allow)
	})
}

func TestService_ResetLoginBucket(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoginBucketRepo = new(mocks.Repository)
		mockIpBucketRepo = new(mocks.Repository)
		mockPasswordBucketRepo = new(mocks.Repository)

		mockLoginBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockIpBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockPasswordBucketRepo.On("DeleteUnusedBucket", mock.Anything)

		key := "login"

		mockLoginBucketRepo.On("ResetBucket", mock.Anything, key).Return(true)
		s := bucket.NewService(mockIpBucketRepo, mockLoginBucketRepo, mockPasswordBucketRepo)

		allow := s.ResetLoginBucket(context.TODO(), key)
		assert.True(t, allow)
	})
}

func TestService_ResetIpBucket(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoginBucketRepo = new(mocks.Repository)
		mockIpBucketRepo = new(mocks.Repository)
		mockPasswordBucketRepo = new(mocks.Repository)

		mockLoginBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockIpBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockPasswordBucketRepo.On("DeleteUnusedBucket", mock.Anything)

		key := "login"

		mockIpBucketRepo.On("ResetBucket", mock.Anything, key).Return(true)
		s := bucket.NewService(mockIpBucketRepo, mockLoginBucketRepo, mockPasswordBucketRepo)

		allow := s.ResetIpBucket(context.TODO(), key)
		assert.True(t, allow)
	})
}
