package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/MaxFando/rate-limiter/internal/config"
	"github.com/MaxFando/rate-limiter/internal/delivery/http/v1/auth"
	"github.com/MaxFando/rate-limiter/internal/providers"
	"github.com/MaxFando/rate-limiter/internal/store/postgres"
	authUC "github.com/MaxFando/rate-limiter/internal/usecase/auth"
	"github.com/MaxFando/rate-limiter/pkg/utils"
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
