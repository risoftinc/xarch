package entities

import "github.com/labstack/echo/v4"

type Response struct {
	Status           int         `json:"status"`
	Message          string      `json:"message"`
	Meta             interface{} `json:"meta,omitempty"`
	Data             interface{} `json:"data,omitempty"`
	Error            interface{} `json:"error,omitempty"`
	ValidationErrors interface{} `json:"validation_errors,omitempty"`
}

func ResponseFormater(ctx echo.Context, Status int, data map[string]interface{}) (int, Response) {
	return Status, Response{
		Status,
		GetResponseCodeMessage(Status),
		data["meta"],
		data["data"],
		data["error"],
		data["validation_errors"],
	}
}
