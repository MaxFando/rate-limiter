package bucket

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/MaxFando/rate-limiter/internal/service/bucket"
	mocks "github.com/MaxFando/rate-limiter/mocks/service/bucket"
	"github.com/MaxFando/rate-limiter/pkg/utils"
)

var (
	mockIPBucketRepo       *mocks.Repository
	mockLoginBucketRepo    *mocks.Repository
	mockPasswordBucketRepo *mocks.Repository
)

const keyLogin = "login"

func TestMain(m *testing.M) {
	utils.InitializeLogger()

	os.Exit(m.Run())
}

func TestService_TryGetPermissionInLoginBucket(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoginBucketRepo = new(mocks.Repository)
		mockIPBucketRepo = new(mocks.Repository)
		mockPasswordBucketRepo = new(mocks.Repository)

		limit := 10
		mockLoginBucketRepo.On("TryGetPermissionInBucket", mock.Anything, keyLogin, limit).Return(true)
		s := bucket.NewService(mockIPBucketRepo, mockLoginBucketRepo, mockPasswordBucketRepo)

		allow := s.TryGetPermissionInLoginBucket(context.TODO(), keyLogin, limit)
		assert.True(t, allow)
	})
}

func TestService_TryGetPermissionInPasswordBucket(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoginBucketRepo = new(mocks.Repository)
		mockIPBucketRepo = new(mocks.Repository)
		mockPasswordBucketRepo = new(mocks.Repository)

		mockLoginBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockIPBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockPasswordBucketRepo.On("DeleteUnusedBucket", mock.Anything)

		key := "password"
		limit := 10
		mockPasswordBucketRepo.On("TryGetPermissionInBucket", mock.Anything, key, limit).Return(true)
		s := bucket.NewService(mockIPBucketRepo, mockPasswordBucketRepo, mockPasswordBucketRepo)

		allow := s.TryGetPermissionInPasswordBucket(context.TODO(), key, limit)
		assert.True(t, allow)
	})
}

func TestService_ResetLoginBucket(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoginBucketRepo = new(mocks.Repository)
		mockIPBucketRepo = new(mocks.Repository)
		mockPasswordBucketRepo = new(mocks.Repository)

		mockLoginBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockIPBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockPasswordBucketRepo.On("DeleteUnusedBucket", mock.Anything)

		mockLoginBucketRepo.On("ResetBucket", mock.Anything, keyLogin).Return(true)
		s := bucket.NewService(mockIPBucketRepo, mockLoginBucketRepo, mockPasswordBucketRepo)

		allow := s.ResetLoginBucket(context.TODO(), keyLogin)
		assert.True(t, allow)
	})
}

func TestService_ResetIpBucket(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLoginBucketRepo = new(mocks.Repository)
		mockIPBucketRepo = new(mocks.Repository)
		mockPasswordBucketRepo = new(mocks.Repository)

		mockLoginBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockIPBucketRepo.On("DeleteUnusedBucket", mock.Anything)
		mockPasswordBucketRepo.On("DeleteUnusedBucket", mock.Anything)

		mockIPBucketRepo.On("ResetBucket", mock.Anything, keyLogin).Return(true)
		s := bucket.NewService(mockIPBucketRepo, mockLoginBucketRepo, mockPasswordBucketRepo)

		allow := s.ResetIPBucket(context.TODO(), keyLogin)
		assert.True(t, allow)
	})
}
