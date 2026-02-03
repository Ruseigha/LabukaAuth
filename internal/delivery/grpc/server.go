package grpc

import (
	"fmt"
	"net"

	"github.com/Ruseigha/LabukaAuth/internal/delivery/grpc/handler"
	"github.com/Ruseigha/LabukaAuth/internal/delivery/grpc/interceptor"
	"github.com/Ruseigha/LabukaAuth/internal/delivery/grpc/proto/proto"
	"github.com/Ruseigha/LabukaAuth/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// SetupServer creates and configures gRPC server
func SetupServer(authService usecase.AuthUseCase) *grpc.Server {
	// Create server with interceptors (middleware)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Recovery(), // First: catch panics
			interceptor.Logger(),   // Second: log requests
		),
	)

	// Register auth service
	authHandler := handler.NewAuthHandler(authService)
	proto.RegisterAuthServiceServer(server, authHandler)

	// Register reflection for grpcurl and other tools
	reflection.Register(server)

	return server
}

// StartServer starts gRPC server on specified port
func StartServer(server *grpc.Server, port string) error {
	// Listen on TCP port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Start serving
	return server.Serve(listener)
}
