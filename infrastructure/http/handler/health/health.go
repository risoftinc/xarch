package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/risoftinc/gologger"
	healthServices "github.com/risoftinc/xarch/domain/services/health"
	"github.com/risoftinc/xarch/infrastructure/http/entities"
)

type (
	IHealthHandler interface {
		Metric(ctx echo.Context) error
	}
	HealthHandler struct {
		logger         gologger.Logger
		healthServices healthServices.IHealthServices
	}
)

func NewHealthHandlers(logger gologger.Logger, healthServices healthServices.IHealthServices) IHealthHandler {
	return &HealthHandler{
		logger:         logger,
		healthServices: healthServices,
	}
}

func (handler HealthHandler) Metric(ctx echo.Context) error {
	ctxReq := ctx.Request().Context()

	handler.logger.WithContext(ctxReq).Info("Healt Check Request started").Send()

	metric, err := handler.healthServices.HealthMetric(ctxReq)
	if err != nil {
		return ctx.JSON(entities.ResponseFormater(ctx, http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		}))
	}

	return ctx.JSON(entities.ResponseFormater(ctx, http.StatusOK, map[string]interface{}{
		"data": metric,
	}))
}
