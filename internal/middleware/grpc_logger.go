package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryLoggingInterceptor creates a gRPC unary-server interceptor that logs detailed information
// about each incoming RPC request and corresponding response via the provided Zap logger.
//
// This interceptor logs the following details:
//   - The full RPC method name (e.g., "/package.Service/Method").
//   - Metadata received from the client.
//   - The request message payload.
//   - The response message payload.
//   - The gRPC status code resulting from RPC handling.
//   - The total latency taken to process the request.
//
// This structured logging is crucial for effective debugging, monitoring, tracing, and auditing,
// offering clear insights into RPC behavior and facilitating issue resolution and performance tracking.
//
// Parameters:
//   - logger: A Zap logger instance used to output the structured logs.
//
// Returns:
//   - grpc.UnaryServerInterceptor: A configured interceptor instance ready to be registered with a gRPC server.
func UnaryLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()

		// Extract the incoming metadata (headers) from the context, if available.
		md, _ := metadata.FromIncomingContext(ctx)

		logger.Info("gRPC request started",
			zap.String("method", info.FullMethod),
			zap.Any("metadata", md),
			zap.Any("request", req),
		)

		// Invoke the actual RPC handler method with the provided context and request
		resp, err := handler(ctx, req)

		// Obtain detailed gRPC status information from the error, if present.
		s, _ := status.FromError(err)

		// Log the outcome of the gRPC call, including response payload and total request processing duration.
		logger.Info("gRPC request completed",
			zap.String("method", info.FullMethod),
			zap.Duration("latency", time.Since(start)),
			zap.Int32("grpc_code", int32(s.Code())),
			zap.Any("response", resp),
			zap.Error(err),
		)

		return resp, err
	}
}
