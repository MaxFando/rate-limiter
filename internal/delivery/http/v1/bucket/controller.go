package bucket

import (
	"github.com/MaxFando/rate-limiter/internal/usecase/bucket"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	uc *bucket.UseCase
}

func NewController(uc *bucket.UseCase) *Controller {
	return &Controller{uc: uc}
}

type resetBucketRequest struct {
	Login string `json:"login"`
	Ip    string `json:"ip"`
}

func (ctr *Controller) ResetBucket(c echo.Context) error {
	ctx := c.Request().Context()

	request := new(resetBucketRequest)
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return echo.NewHTTPError(422, err.Error())
	}

	_, err := ctr.uc.Reset(ctx, request.Login, request.Ip)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSON(200, map[string]interface{}{"status": "ok"})
}