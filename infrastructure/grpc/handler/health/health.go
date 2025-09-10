package health

import (
	"context"

	"github.com/risoftinc/gologger"
	healthServices "github.com/risoftinc/xarch/domain/services/health"
	"github.com/risoftinc/xarch/infrastructure/grpc/entities"
	healthpb "github.com/risoftinc/xarch/infrastructure/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	HealthHandler struct {
		healthpb.UnimplementedHealthServiceServer
		logger         gologger.Logger
		healthServices healthServices.IHealthServices
	}
)

func NewHealthHandlers(logger gologger.Logger, healthServices healthServices.IHealthServices) *HealthHandler {
	return &HealthHandler{
		logger:         logger,
		healthServices: healthServices,
	}
}

func (handler HealthHandler) GetHealthMetric(ctx context.Context, req *healthpb.HealthMetricRequest) (*healthpb.HealthMetricResponse, error) {
	// Get request ID and language from context (set by middleware)

	handler.logger.WithContext(ctx).Info("Health Check Request started").Send()

	metric, err := handler.healthServices.HealthMetric(ctx)
	if err != nil {
		handler.logger.WithContext(ctx).Error("Health check failed: " + err.Error()).Send()

		// Create error response
		response := &healthpb.HealthMetricResponse{
			Status:  14,
			Message: entities.GetResponseCodeMessage(500),
			Error:   err.Error(),
		}

		return response, status.Errorf(codes.Internal, "Health check failed: %v", err)
	}

	// Convert health metric to protobuf format
	statusMap := make(map[string]string)
	for k, v := range metric.Status {
		if str, ok := v.(string); ok {
			statusMap[k] = str
		} else {
			statusMap[k] = "unknown"
		}
	}

	// Create database info

	response := &healthpb.HealthMetricResponse{
		Status:  0,
		Message: entities.GetResponseCodeMessage(200),
		Data: &healthpb.HealthMetricData{
			Status: statusMap,
			Database: &healthpb.DatabaseInfo{
				MaxOpenConnections: int32(metric.DB.MaxOpenConnections),
				OpenConnections:    int32(metric.DB.OpenConnections),
				InUse:              int32(metric.DB.InUse),
				Idle:               int32(metric.DB.Idle),
				WaitCount:          int32(metric.DB.WaitCount),
				WaitDuration:       int32(metric.DB.WaitDuration),
				MaxIdleClosed:      int32(metric.DB.MaxIdleClosed),
				MaxIdleTimeClosed:  int32(metric.DB.MaxIdleTimeClosed),
				MaxLifetimeClosed:  int32(metric.DB.MaxLifetimeClosed),
			},
		},
	}

	handler.logger.WithContext(ctx).Info("Health Check Request completed successfully").Send()
	return response, nil
}
