package auth

import (
	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/usecase/auth"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	uc *auth.UseCase
}

func NewAuthController(usecase *auth.UseCase) *Controller {
	return &Controller{
		uc: usecase,
	}
}

type tryAuthorizationRequest struct {
	Login    string `json:"login" query:"login" validate:"required"`
	Password string `json:"password" query:"password" validate:"required"`
	Ip       string `json:"ip" query:"ip" validate:"required"`
}

func (ctr *Controller) TryAuthorization(c echo.Context) error {
	ctx := c.Request().Context()

	request := new(tryAuthorizationRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(500, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(422, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	payload := network.Request{Login: request.Login, Password: request.Password, Ip: request.Ip}
	allow, err := ctr.uc.TryAuthorization(ctx, payload)
	if err != nil {
		return c.JSON(500, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	return c.JSON(200, map[string]interface{}{"ok": allow})
}
