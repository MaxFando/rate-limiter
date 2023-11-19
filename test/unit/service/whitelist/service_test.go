package whitelist

import (
	"context"
	"errors"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/service/whitelist"
	mocks "github.com/MaxFando/rate-limiter/mocks/service/whitelist"
	"github.com/MaxFando/rate-limiter/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

var (
	mockStore *mocks.Store
)

func TestMain(m *testing.M) {
	utils.InitializeLogger()

	os.Exit(m.Run())
}

func TestService_AddIP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStore = new(mocks.Store)
		s := whitelist.NewService(mockStore)

		payload := network.IpNetwork{
			Ip:   "192.168.1.1",
			Mask: "255.255.255.0",
		}
		prefix, _ := utils.GetPrefix(payload.Ip, payload.Mask)

		mockStore.On("Add", mock.Anything, prefix, payload.Mask).Return(nil)

		err := s.AddIP(context.TODO(), payload)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockStore = new(mocks.Store)
		s := whitelist.NewService(mockStore)

		payload := network.IpNetwork{
			Ip:   "192.168.1.1",
			Mask: "255.255.255.0",
		}
		prefix, _ := utils.GetPrefix(payload.Ip, payload.Mask)

		errMock := errors.New("test")
		mockStore.On("Add", mock.Anything, prefix, payload.Mask).Return(errMock)

		err := s.AddIP(context.TODO(), payload)
		assert.ErrorAs(t, err, &errMock)
	})
}
