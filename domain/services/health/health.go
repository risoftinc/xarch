package health

import (
	"context"
	"strings"

	"go.risoftinc.com/gologger"
	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/constant"
	healthModels "go.risoftinc.com/xarch/domain/models/health"
	healthRepositories "go.risoftinc.com/xarch/domain/repositories/health"
)

type (
	IHealthServices interface {
		HealthMetric(ctx context.Context) (*healthModels.HealthMetric, error)
	}
	HealthServices struct {
		logger             gologger.Logger
		healthRepositories healthRepositories.IHealthRepositories
	}
)

func NewHealthService(
	logger gologger.Logger,
	healthRepositories healthRepositories.IHealthRepositories,
) IHealthServices {
	return &HealthServices{
		logger:             logger,
		healthRepositories: healthRepositories,
	}
}

func (svc HealthServices) HealthMetric(ctx context.Context) (*healthModels.HealthMetric, error) {
	metric := &healthModels.HealthMetric{
		Status: make(map[string]interface{}),
	}

	// Check database health
	databaseHealth, err := svc.healthRepositories.DatabaseHealth(ctx)
	if err != nil {
		metric.Status["database"] = "disconnected"
		svc.logger.WithContext(ctx).Error("Error database health").ErrorData(err).Send()
		return metric, goresponse.NewResponseBuilder(categorizeError(err)).
			WithContext(ctx).
			SetError(err).
			ToError()
	}

	metric.Status["database"] = "connected"
	metric.DB = databaseHealth

	return metric, nil
}

func categorizeError(err error) string {
	errMsg := strings.ToLower(err.Error())

	// Connection errors -> 503/UNAVAILABLE
	if strings.Contains(errMsg, "connection refused") {
		return constant.ErrorConnectionRefused
	}
	if strings.Contains(errMsg, "too many connections") {
		return constant.ErrorTooManyConnections
	}
	if strings.Contains(errMsg, "timeout") {
		return constant.ErrorConnectionTimeout
	}
	if strings.Contains(errMsg, "no such host") {
		return constant.ErrorDnsError
	}

	// Auth errors -> 401/UNAUTHENTICATED
	if strings.Contains(errMsg, "authentication failed") {
		return constant.ErrorAuthFailed
	}
	if strings.Contains(errMsg, "access denied") {
		return constant.ErrorAccessDenied
	}

	// Internal errors -> 500/INTERNAL
	if strings.Contains(errMsg, "driver") {
		return constant.ErrorDriverError
	}
	if strings.Contains(errMsg, "ssl") || strings.Contains(errMsg, "tls") {
		return constant.ErrorSslTlsError
	}

	// Default -> 500/INTERNAL
	return constant.ErrorInternalServer
}
