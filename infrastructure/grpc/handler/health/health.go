package health

import (
	"context"

	"go.risoftinc.com/gologger"
	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/constant"
	healthServices "go.risoftinc.com/xarch/domain/services/health"
	"go.risoftinc.com/xarch/infrastructure/grpc/entities"
	healthpb "go.risoftinc.com/xarch/infrastructure/grpc/proto"
	"go.risoftinc.com/xarch/utils/grpc"
	"google.golang.org/grpc/status"
)

type (
	HealthHandler struct {
		healthpb.UnimplementedHealthServiceServer
		logger         gologger.Logger
		grpcEntities   entities.IGrpcEntities
		healthServices healthServices.IHealthServices
	}
)

func NewHealthHandlers(
	logger gologger.Logger,
	grpcEntities entities.IGrpcEntities,
	healthServices healthServices.IHealthServices,
) *HealthHandler {
	return &HealthHandler{
		logger:         logger,
		grpcEntities:   grpcEntities,
		healthServices: healthServices,
	}
}

func (handler HealthHandler) GetHealthMetric(ctx context.Context, req *healthpb.HealthMetricRequest) (*healthpb.HealthMetricResponse, error) {
	handler.logger.WithContext(ctx).Info("Health Check Request started").Send()

	metric, err := handler.healthServices.HealthMetric(ctx)
	if err != nil {
		handler.logger.WithContext(ctx).Error("Health check failed: " + err.Error()).Send()

		// Use gRPC response formatter for error
		grpcResponse := handler.grpcEntities.ResponseFormaterError(err)

		return &healthpb.HealthMetricResponse{
			Meta: &healthpb.Meta{
				Message: grpcResponse.Meta.Message,
				Error:   &grpcResponse.Meta.Error, // Will be nil if no error, so field won't appear in JSON
			},
		}, status.Errorf(grpc.IntToCode(grpcResponse.Code), "%s", grpcResponse.Meta.Message)
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

	// Use goresponse for success response
	responseBuilder := goresponse.NewResponseBuilder(constant.IsResponseSuccess).
		WithContext(ctx).SetData("data", metric)

	// Use gRPC response formatter for success
	grpcResponse := handler.grpcEntities.ResponseFormater(responseBuilder)

	// Only set error if it's not empty (for success case, error should be nil)
	var errorPtr *string
	if grpcResponse.Meta.Error != "" {
		errorPtr = &grpcResponse.Meta.Error
	}

	// Convert to protobuf response with new structure
	response := &healthpb.HealthMetricResponse{
		Meta: &healthpb.Meta{
			Message: grpcResponse.Meta.Message,
			Error:   errorPtr, // Will be nil for success, so field won't appear in JSON
		},
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
