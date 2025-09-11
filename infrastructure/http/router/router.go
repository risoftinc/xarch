package router

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	dep "go.risoftinc.com/xarch/infrastructure/http"
	"go.risoftinc.com/xarch/utils/validator"
)

func Routers(dep *dep.Dependencies) *echo.Echo {
	engine := echo.New()

	// Add custom validator
	engine.Validator = validator.NewCustomValidator()

	// Add request ID middleware globally
	engine.Use(dep.Middlewares.ContextMiddleware())
	engine.Use(echoMiddleware.Recover())

	// Public routes
	engine.GET("/health", dep.HealthHandlers.Metric)

	return engine
}
