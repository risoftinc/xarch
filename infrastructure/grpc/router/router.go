package router

import (
	dep "github.com/risoftinc/xarch/infrastructure/grpc"
	healthpb "github.com/risoftinc/xarch/infrastructure/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RegisterGRPCServices registers all gRPC services
func RegisterGRPCServices(dep *dep.Dependencies) *grpc.Server {
	// Initialize gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(dep.Middlewares.UnaryContextInterceptor()),
		grpc.StreamInterceptor(dep.Middlewares.StreamContextInterceptor()),
	)

	// Register health service
	healthpb.RegisterHealthServiceServer(grpcServer, dep.HealthHandlers)

	// Enable reflection for debugging and testing
	reflection.Register(grpcServer)

	return grpcServer
}
