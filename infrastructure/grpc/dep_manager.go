//go:build elsabuild
// +build elsabuild

package http

import (
	"github.com/risoftinc/elsa"
	"github.com/risoftinc/gologger"
	"github.com/risoftinc/xarch/config"
	healthRepo "github.com/risoftinc/xarch/domain/repositories/health"
	healthSvc "github.com/risoftinc/xarch/domain/services/health"
	healthHandler "github.com/risoftinc/xarch/infrastructure/grpc/handler/health"
	mid "github.com/risoftinc/xarch/infrastructure/grpc/middleware"
	"gorm.io/gorm"
)

type Dependencies struct {
	Middlewares    mid.IContextMiddleware
	HealthHandlers healthHandler.HealthHandler
}

func InitializeServices(db *gorm.DB, cfg config.Config, logger gologger.Logger) *Dependencies {
	elsa.Generate(
		RepositorySet,
		ServicesSet,
		HandlerSet,
		MidlewareSet,
	)

	return nil
}

var RepositorySet = elsa.Set(
	healthRepo.NewHealthRepositories,
)

var ServicesSet = elsa.Set(
	healthSvc.NewHealthService,
)

var HandlerSet = elsa.Set(
	healthHandler.NewHealthHandlers,
)

var MidlewareSet = elsa.Set(
	mid.NewContextMiddleware,
)
