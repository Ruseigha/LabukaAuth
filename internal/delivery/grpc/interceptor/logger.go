package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// Logger logs gRPC requests
func Logger() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Call handler
		resp, err := handler(ctx, req)

		// Log
		duration := time.Since(start)
		if err != nil {
			log.Printf("gRPC %s %s ERROR: %v", info.FullMethod, duration, err)
		} else {
			log.Printf("gRPC %s %s OK", info.FullMethod, duration)
		}

		return resp, err
	}
}
