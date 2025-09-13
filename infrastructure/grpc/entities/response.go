package entities

import (
	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/constant"
	"go.risoftinc.com/xarch/domain/models/response"
)

// IGrpcEntities interface for gRPC response formatting
type (
	IGrpcEntities interface {
		ResponseFormaterError(err error) *response.Response
		ResponseFormater(res *goresponse.ResponseBuilder) *response.Response
	}

	// GrpcEntities struct for gRPC response handling
	GrpcEntities struct {
		async    *goresponse.AsyncConfigManager
		response *goresponse.ResponseConfig
	}
)

// NewGrpcEntities creates a new gRPC entities instance
func NewGrpcEntities(async *goresponse.AsyncConfigManager) IGrpcEntities {
	e := &GrpcEntities{
		async:    async,
		response: async.GetConfig(),
	}

	async.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		e.response = newConfig
	})

	return e
}

// ResponseFormaterError handles error responses for gRPC
func (e *GrpcEntities) ResponseFormaterError(err error) *response.Response {
	var rb *goresponse.ResponseBuilder
	res, ok := goresponse.ParseResponseBuilderError(err)
	if !ok {
		rb = goresponse.NewResponseBuilder(constant.ErrorInternalServer).SetError(err)
	} else {
		rb = res
	}

	resBuild, err := e.response.BuildResponse(rb)
	if err != nil {
		return &response.Response{
			Code: 13,
			Meta: response.Meta{
				Message: "Internal Server Error",
				Error:   err.Error(),
			},
		}
	}

	return &response.Response{
		Code: resBuild.Code,
		Meta: response.Meta{
			Message: resBuild.Message,
			Error:   resBuild.Error.Error(),
		},
	}
}

// ResponseFormater handles success responses for gRPC
func (e *GrpcEntities) ResponseFormater(res *goresponse.ResponseBuilder) *response.Response {
	resBuild, err := e.response.BuildResponse(res)
	if err != nil {
		return &response.Response{
			Code: 0,
			Meta: response.Meta{
				Message: "OK",
				Error:   err.Error(),
			},
			Data: res.Data,
		}
	}

	return &response.Response{
		Code: resBuild.Code,
		Meta: response.Meta{
			Message: resBuild.Message,
		},
		Data: resBuild.Data["data"],
	}
}
