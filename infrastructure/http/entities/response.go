package entities

import (
	"github.com/labstack/echo/v4"
	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/constant"
	"go.risoftinc.com/xarch/domain/models/response"
)

// {
//     "meta": {
//         "message": "Category created successfully",
//         "error": ""
//     },
//     "data": {
//         "id": 3,
//         "name": "Makanan"
//     }
// }

type (
	IEntities interface {
		ResponseFormaterError(ctx echo.Context, err error) error
		ResponseFormater(ctx echo.Context, res *goresponse.ResponseBuilder) error
	}
	Entities struct {
		async    *goresponse.AsyncConfigManager
		response *goresponse.ResponseConfig
	}
)

func NewEntities(async *goresponse.AsyncConfigManager) IEntities {
	e := &Entities{
		async:    async,
		response: async.GetConfig(),
	}

	async.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		e.response = newConfig
	})

	return e
}

func (e *Entities) ResponseFormaterError(ctx echo.Context, err error) error {
	var rb *goresponse.ResponseBuilder
	res, ok := goresponse.ParseResponseBuilderError(err)
	if !ok {
		rb = goresponse.NewResponseBuilder(constant.ErrorInternalServer).SetError(err)
	} else {
		rb = res
	}

	resBuild, err := e.response.BuildResponse(rb)
	if err != nil {
		return ctx.JSON(500, response.Response{
			Meta: response.Meta{
				Message: "Internal Server Error",
				Error:   err.Error(),
			},
			Data: e.response,
		})
	}

	return ctx.JSON(resBuild.Code, response.Response{
		Meta: response.Meta{
			Message: resBuild.Message,
			Error:   resBuild.Error.Error(),
		},
	})
}

func (e *Entities) ResponseFormater(ctx echo.Context, res *goresponse.ResponseBuilder) error {
	resBuild, err := e.response.BuildResponse(res)
	if err != nil {
		return ctx.JSON(200, response.Response{
			Meta: response.Meta{
				Message: "OK",
				Error:   err.Error(),
			},
			Data: res.Data,
		})
	}

	return ctx.JSON(resBuild.Code, response.Response{
		Meta: response.Meta{
			Message: resBuild.Message,
		},
		Data: resBuild.Data["data"],
	})
}
