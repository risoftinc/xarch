package entities

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Response represents a gRPC response structure
type Response struct {
	Status           int32       `json:"status"`
	Message          string      `json:"message"`
	Meta             interface{} `json:"meta,omitempty"`
	Data             interface{} `json:"data,omitempty"`
	Error            string      `json:"error,omitempty"`
	ValidationErrors interface{} `json:"validation_errors,omitempty"`
}

// ResponseFormater formats the response for gRPC
func ResponseFormater(ctx context.Context, statusCode int32, data map[string]interface{}) (*Response, error) {
	response := &Response{
		Status:  statusCode,
		Message: GetGRPCStatus(statusCode).String(),
		Meta:    data["meta"],
		Data:    data["data"],
		Error:   "",
	}

	if data["error"] != nil {
		if errStr, ok := data["error"].(string); ok {
			response.Error = errStr
		}
	}

	if data["validation_errors"] != nil {
		response.ValidationErrors = data["validation_errors"]
	}

	return response, nil
}

// GetGRPCStatus converts HTTP status to gRPC status
func GetGRPCStatus(httpStatus int32) codes.Code {
	switch {
	case httpStatus >= 200 && httpStatus < 300:
		return codes.OK
	case httpStatus == 400:
		return codes.InvalidArgument
	case httpStatus == 401:
		return codes.Unauthenticated
	case httpStatus == 403:
		return codes.PermissionDenied
	case httpStatus == 404:
		return codes.NotFound
	case httpStatus == 409:
		return codes.AlreadyExists
	case httpStatus == 422:
		return codes.InvalidArgument
	case httpStatus >= 500:
		return codes.Internal
	default:
		return codes.Unknown
	}
}

// CreateGRPCError creates a gRPC error from HTTP status and message
func CreateGRPCError(httpStatus int32, message string) error {
	grpcCode := GetGRPCStatus(httpStatus)
	return status.Errorf(grpcCode, message)
}
