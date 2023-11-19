package whitelist

import (
	"github.com/labstack/echo/v4"

	"github.com/MaxFando/rate-limiter/internal/domain/network"
	"github.com/MaxFando/rate-limiter/internal/usecase/whitelist"
)

type Controller struct {
	uc *whitelist.UseCase
}

func NewController(usecase *whitelist.UseCase) *Controller {
	return &Controller{
		uc: usecase,
	}
}

type addIpRequest struct {
	Ip   string `json:"ip" query:"ip" validate:"required"`
	Mask string `json:"mask" query:"mask" validate:"required"`
}

func (ctr *Controller) AddIP(c echo.Context) error {
	ctx := c.Request().Context()

	request := new(addIpRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(500, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(422, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	payload, err := network.NewIpNetwork(request.Ip, request.Mask)
	if err != nil {
		return c.JSON(500, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	err = ctr.uc.AddIP(ctx, payload)
	if err != nil {
		return c.JSON(500, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	return c.JSON(200, map[string]interface{}{"ok": true})
}

type removeIpRequest struct {
	Ip   string `json:"ip" query:"ip" validate:"required"`
	Mask string `json:"mask" query:"mask" validate:"required"`
}

func (ctr *Controller) RemoveIP(c echo.Context) error {
	ctx := c.Request().Context()

	request := new(removeIpRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(500, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(422, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	payload, err := network.NewIpNetwork(request.Ip, request.Mask)
	if err != nil {
		return c.JSON(500, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	err = ctr.uc.RemoveIP(ctx, payload)
	if err != nil {
		return c.JSON(500, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	return c.JSON(200, map[string]interface{}{"ok": true})
}

func (ctr *Controller) GetIPList(c echo.Context) error {
	ctx := c.Request().Context()

	ipList, err := ctr.uc.GetIPList(ctx)
	if err != nil {
		return c.JSON(500, map[string]interface{}{"ok": false, "error": err.Error()})
	}

	return c.JSON(200, map[string]interface{}{"ok": true, "data": ipList})
}
