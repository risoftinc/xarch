package health

import (
	"context"
	"strings"

	"go.risoftinc.com/gologger"
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

		return metric, err
	} else {
		metric.Status["database"] = "connected"
		metric.DB = databaseHealth
	}

	return metric, nil
}

func categorizeError(err error) (string, int, int) {
	errMsg := strings.ToLower(err.Error())

	// Connection errors -> 503/UNAVAILABLE
	if strings.Contains(errMsg, "connection refused") {
		return "CONNECTION_REFUSED", 503, 14
	}
	if strings.Contains(errMsg, "too many connections") {
		return "TOO_MANY_CONNECTIONS", 503, 14
	}
	if strings.Contains(errMsg, "timeout") {
		return "CONNECTION_TIMEOUT", 503, 14
	}
	if strings.Contains(errMsg, "no such host") {
		return "DNS_ERROR", 503, 14
	}

	// Auth errors -> 401/UNAUTHENTICATED
	if strings.Contains(errMsg, "authentication failed") {
		return "AUTH_FAILED", 401, 16
	}
	if strings.Contains(errMsg, "access denied") {
		return "ACCESS_DENIED", 401, 16
	}

	// Internal errors -> 500/INTERNAL
	if strings.Contains(errMsg, "driver") {
		return "DRIVER_ERROR", 500, 13
	}
	if strings.Contains(errMsg, "ssl") || strings.Contains(errMsg, "tls") {
		return "SSL_TLS_ERROR", 500, 13
	}

	// Default -> 500/INTERNAL
	return "UNKNOWN_ERROR", 500, 13
}
