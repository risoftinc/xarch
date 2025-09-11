//go:build elsabuild
// +build elsabuild

package grpc

import (
	"go.risoftinc.com/elsa"
	"go.risoftinc.com/gologger"
	"go.risoftinc.com/xarch/config"
	healthRepo "go.risoftinc.com/xarch/domain/repositories/health"
	healthSvc "go.risoftinc.com/xarch/domain/services/health"
	healthHandler "go.risoftinc.com/xarch/infrastructure/grpc/handler/health"
	mid "go.risoftinc.com/xarch/infrastructure/grpc/middleware"
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
