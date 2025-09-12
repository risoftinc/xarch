package health

import (
	"github.com/labstack/echo/v4"
	"go.risoftinc.com/gologger"
	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/constant"
	healthServices "go.risoftinc.com/xarch/domain/services/health"
	"go.risoftinc.com/xarch/infrastructure/http/entities"
)

type (
	IHealthHandler interface {
		Metric(ctx echo.Context) error
	}
	HealthHandler struct {
		logger         gologger.Logger
		entities       entities.IEntities
		healthServices healthServices.IHealthServices
	}
)

func NewHealthHandlers(
	logger gologger.Logger,
	entities entities.IEntities,
	healthServices healthServices.IHealthServices,
) IHealthHandler {
	return &HealthHandler{
		logger:         logger,
		entities:       entities,
		healthServices: healthServices,
	}
}

func (handler HealthHandler) Metric(ctx echo.Context) error {
	ctxReq := ctx.Request().Context()

	handler.logger.WithContext(ctxReq).Info("Healt Check Request started").Send()

	metric, err := handler.healthServices.HealthMetric(ctxReq)
	if err != nil {
		return handler.entities.ResponseFormaterError(ctx, err)
	}

	return handler.entities.ResponseFormater(ctx,
		goresponse.NewResponseBuilder(constant.IsResponseSuccess).
			WithContext(ctxReq).SetData("data", metric),
	)
}
