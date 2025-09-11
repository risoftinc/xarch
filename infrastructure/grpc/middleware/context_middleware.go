package middleware

import (
	"context"

	"github.com/google/uuid"
	"go.risoftinc.com/gologger"
	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	RequestIDHeader = "x-request-id"
	LanguageHeader  = "x-language"
)

type (
	IContextMiddleware interface {
		UnaryContextInterceptor() grpc.UnaryServerInterceptor
		StreamContextInterceptor() grpc.StreamServerInterceptor
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

// UnaryContextInterceptor adds request context data (request ID, language, etc.) to unary gRPC calls
func (cm ContextMiddleware) UnaryContextInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		// Get request ID and language from metadata
		requestID := getMetadataValue(md, RequestIDHeader)
		language := getMetadataValue(md, LanguageHeader)

		// If no request ID in metadata, generate a new UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// If no language in metadata, set default language
		if language == "" {
			language = constant.DefaultLanguage
		}

		// Create new context with request ID and language
		ctx = gologger.WithRequestID(ctx, requestID)
		ctx = goresponse.WithLanguage(ctx, language)
		ctx = goresponse.WithProtocol(ctx, constant.ProtocolGrpc)

		// Add metadata to outgoing context for client reference
		ctx = metadata.AppendToOutgoingContext(ctx,
			RequestIDHeader, requestID,
			LanguageHeader, language,
		)

		// Log the incoming request
		cm.logger.WithContext(ctx).Info("gRPC request started").
			Data("method", info.FullMethod).
			Data("request_id", requestID).
			Data("language", language).
			Send()

		// Call the actual handler
		resp, err := handler(ctx, req)

		// Log the response
		if err != nil {
			cm.logger.WithContext(ctx).Error("gRPC request failed").
				Data("method", info.FullMethod).
				Data("error", err.Error()).
				Send()
		} else {
			cm.logger.WithContext(ctx).Info("gRPC request completed").
				Data("method", info.FullMethod).
				Send()
		}

		return resp, err
	}
}

// StreamContextInterceptor adds request context data (request ID, language, etc.) to streaming gRPC calls
func (cm ContextMiddleware) StreamContextInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Extract metadata from context
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		// Get request ID and language from metadata
		requestID := getMetadataValue(md, RequestIDHeader)
		language := getMetadataValue(md, LanguageHeader)

		// If no request ID in metadata, generate a new UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// If no language in metadata, set default language
		if language == "" {
			language = constant.DefaultLanguage
		}

		// Create new context with request ID and language
		ctx = gologger.WithRequestID(ctx, requestID)
		ctx = goresponse.WithLanguage(ctx, language)
		ctx = goresponse.WithProtocol(ctx, constant.ProtocolGrpc)

		// Add metadata to outgoing context for client reference
		ctx = metadata.AppendToOutgoingContext(ctx,
			RequestIDHeader, requestID,
			LanguageHeader, language,
		)

		// Create a new server stream with the updated context
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// Log the incoming stream request
		cm.logger.WithContext(ctx).Info("gRPC stream request started").
			Data("method", info.FullMethod).
			Data("request_id", requestID).
			Data("language", language).
			Send()

		// Call the actual handler
		err := handler(srv, wrappedStream)

		// Log the stream response
		if err != nil {
			cm.logger.WithContext(ctx).Error("gRPC stream request failed").
				Data("method", info.FullMethod).
				Data("error", err.Error()).
				Send()
		} else {
			cm.logger.WithContext(ctx).Info("gRPC stream request completed").
				Data("method", info.FullMethod).
				Send()
		}

		return err
	}
}

// wrappedServerStream wraps the original ServerStream with a new context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// getMetadataValue safely extracts a value from gRPC metadata
func getMetadataValue(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// GetRequestIDFromContext extracts request ID from context
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(gologger.RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetLanguageFromContext extracts language from context
func GetLanguageFromContext(ctx context.Context) string {
	if language, ok := ctx.Value(goresponse.LanguageKey).(string); ok {
		return language
	}
	return constant.DefaultLanguage
}
