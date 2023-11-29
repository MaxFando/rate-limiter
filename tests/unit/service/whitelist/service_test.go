package whitelist

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/service/whitelist"
	mocks "github.com/MaxFando/rate-limiter/mocks/service/whitelist"
	"github.com/MaxFando/rate-limiter/pkg/utils"
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

		payload, _ := network.NewIPNetwork(
			"192.168.1.1",
			"255.255.255.0",
		)
		prefix, _ := utils.GetPrefix(payload.IP.String(), payload.Mask.String())

		mockStore.On("Add", mock.Anything, prefix, payload.Mask.String()).Return(nil)

		err := s.AddIP(context.TODO(), payload)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockStore = new(mocks.Store)
		s := whitelist.NewService(mockStore)

		payload, _ := network.NewIPNetwork(
			"192.168.1.1",
			"255.255.255.0",
		)
		prefix, _ := utils.GetPrefix(payload.IP.String(), payload.Mask.String())

		errMock := errors.New("tests")
		mockStore.On("Add", mock.Anything, prefix, payload.Mask.String()).Return(errMock)

		err := s.AddIP(context.TODO(), payload)
		assert.ErrorAs(t, err, &errMock)
	})
}
