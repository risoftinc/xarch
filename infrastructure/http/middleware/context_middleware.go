package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.risoftinc.com/gologger"
	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/constant"
)

const (
	RequestIDHeader = "X-Request-ID"
	LanguageHeader  = "X-Language"
)

type (
	IContextMiddleware interface {
		ContextMiddleware() echo.MiddlewareFunc
	}
	ContextMiddleware struct {
		logger gologger.Logger
	}
)

func NewContextMiddleware(logger gologger.Logger) IContextMiddleware {
	return &ContextMiddleware{
		logger: logger,
	}
}

// ContextMiddleware adds request context data (request ID, language, etc.) to the context
func (rm ContextMiddleware) ContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if data exists in header
			requestID := c.Request().Header.Get(RequestIDHeader)
			language := c.Request().Header.Get(LanguageHeader)

			// If no request ID in header, generate a new UUID
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// If no language in header, set default language
			if language == "" {
				language = constant.DefaultLanguage
			}

			// Set data to context echo
			c.Set(string(gologger.RequestIDKey), requestID)
			c.Set(string(goresponse.LanguageKey), language)

			// Set context
			ctx := gologger.WithRequestID(c.Request().Context(), requestID)
			ctx = goresponse.WithLanguage(ctx, language)
			ctx = goresponse.WithProtocol(ctx, constant.ProtocolWebApi)

			// Set data to context request
			c.SetRequest(c.Request().WithContext(ctx))

			// Add data to response header for client reference
			c.Response().Header().Set(RequestIDHeader, requestID)
			c.Response().Header().Set(LanguageHeader, language)

			return next(c)
		}
	}
}
