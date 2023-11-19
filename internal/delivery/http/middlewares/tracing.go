package middlewares

import (
	"context"
	"fmt"
	"github.com/MaxFando/rate-limiter/pkg/tracing"
	"github.com/MaxFando/rate-limiter/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"time"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := tracing.CreateChildSpan(c.Request().Context(), fmt.Sprintf("HTTP GET URL: %s", c.Path()))
		defer span.Finish()

		var traceId string
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			traceId = sc.TraceID().String()
		}

		c.Response().Header().Set("X-Trace-Id", traceId)
		c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "trace_id", traceId)))
		c.Response().Header().Set("X-Trace-Id", traceId)

		utils.Logger.Info(
			"Server request",
			zap.String("method", c.Request().Method),
			zap.String("URL", c.Request().RequestURI),
		)

		timeStart := time.Now()

		if err := next(c); err != nil {
			c.Error(err)
		}

		statusOK := c.Response().Status >= 200 && c.Response().Status < 300
		if !statusOK {
			span.LogKV("error", true)
		}

		timeEnd := time.Now()
		diffTime := timeEnd.Sub(timeStart)

		utils.Logger.Info(
			"Server response",
			zap.Int64("response length", c.Response().Size),
			zap.Duration("time", diffTime),
		)

		return nil
	}
}
