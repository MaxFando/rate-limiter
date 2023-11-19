package auth

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/config"
	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/auth"
	"github.com/MaxFando/rate-limiter/internal/providers"
	"github.com/MaxFando/rate-limiter/internal/store/postgres"
	authUC "github.com/MaxFando/rate-limiter/internal/usecase/auth"
	"github.com/MaxFando/rate-limiter/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var authControllerV1 *auth.Controller

func TestMain(m *testing.M) {
	utils.InitializeLogger()

	ctx := context.TODO()
	config.InitializeConfig()

	postgresDB := postgres.New()
	_ = postgresDB.Connect(ctx, "default", config.Config.Database)

	repositoryProvider := providers.NewRepositoryProvider(postgresDB)
	repositoryProvider.BootPrefixProviderContract()

	serviceProvider := providers.NewServiceProvider()
	serviceProvider.RegisterDependencies(repositoryProvider)

	authUseCase := authUC.NewUseCase(serviceProvider.BlacklistService, serviceProvider.WhitelistService, serviceProvider.BucketService)
	authControllerV1 = auth.NewAuthController(authUseCase)

	os.Exit(m.Run())
}

var testJSON = `{"login":"login", "password": "password", "ip": "192.168.1.1"}`

func TestController_TryAuthorization(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(testJSON))
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	err := authControllerV1.TryAuthorization(c)
	assert.NoError(t, err)
}

//func TestTryAuthorization(t *testing.T) {
//	logger := zap.NewExample().Sugar()
//
//	controller := gomock.NewController(t)
//	defer controller.Finish()
//
//	blackListMockStore := mock_service.NewMockBlackListStore(controller)
//	blackList := service.NewBlackList(blackListMockStore, logger)
//
//	whiteListMockStore := mock_service.NewMockWhiteListStore(controller)
//	whiteList := service.NewWhiteList(whiteListMockStore, logger)
//
//	cfg, err := config.New()
//	if err != nil {
//		require.NoError(t, err)
//	}
//
//	serviceAuth := service.NewAuthorization(blackList, whiteList, cfg, logger)
//
//	auth := NewAuthorization(serviceAuth, logger)
//
//	cases := []struct {
//		name    string
//		request entity.Request
//	}{
//		{name: "valid request", request: entity.Request{
//			Login:    "test",
//			Password: "1234",
//			Ip:       "192.1.5.1",
//		}},
//	}
//
//	blackListMockStore.EXPECT().Get().Return([]entity.IpNetwork{}, nil).AnyTimes()
//	whiteListMockStore.EXPECT().Get().Return([]entity.IpNetwork{}, nil).AnyTimes()
//
//	router := httprouter.New()
//	router.POST("/auth/check", auth.TryAuthorization)
//
//	request := cases[0].request
//
//	body, err := json.Marshal(request)
//	require.NoError(t, err)
//
//	req, err := http.NewRequest("POST", "/auth/check", bytes.NewReader(body))
//	require.NoError(t, err)
//	req.Header.Set("Content-Type", "application/json")
//	rr := httptest.NewRecorder()
//	router.ServeHTTP(rr, req)
//
//	require.Equal(t, http.StatusOK, rr.Code)
//	s := rr.Body.String()
//	require.Equal(t, "ok=true", s)
//
//}
