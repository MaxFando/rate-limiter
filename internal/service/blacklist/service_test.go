package blacklist

import (
	"context"
	"errors"
	"testing"

	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_AddIP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockStore(ctrl)
	service := NewService(mockStore)

	ipNetwork, err := network.NewIPNetwork("192.168.1.1", "255.255.255.0")
	assert.NoError(t, err)
	mockStore.EXPECT().Add(gomock.Any(), "192.168.1.0", "255.255.255.0").Return(nil)

	err = service.AddIP(context.Background(), ipNetwork)
	assert.NoError(t, err)
}

func TestService_AddIP_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockStore(ctrl)
	service := NewService(mockStore)

	ipNetwork, err := network.NewIPNetwork("192.168.1.1", "255.255.255.0")
	assert.NoError(t, err)

	mockStore.EXPECT().Add(gomock.Any(), "192.168.1.0", "255.255.255.0").Return(errors.New("add error"))

	err = service.AddIP(context.Background(), ipNetwork)
	assert.Error(t, err)
}

func TestService_RemoveIP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer ctrl.Finish()

	mockStore := NewMockStore(ctrl)
	service := NewService(mockStore)

	ipNetwork, err := network.NewIPNetwork("192.168.1.1", "255.255.255.0")
	assert.NoError(t, err)

	mockStore.EXPECT().Remove(gomock.Any(), "192.168.1.0", "255.255.255.0").Return(errors.New("remove error"))

	err = service.RemoveIP(context.Background(), ipNetwork)
	assert.Error(t, err)
}
