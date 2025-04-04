package middleware

import (
    "context"
    "time"

    "go.uber.org/zap"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
)

// UnaryLoggingInterceptor logs each unary RPC request and response
func UnaryLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {

        start := time.Now()

        // Retrieve metadata if needed
        md, _ := metadata.FromIncomingContext(ctx)

        logger.Info("gRPC request start",
            zap.String("method", info.FullMethod),
            zap.Any("metadata", md),
            zap.Any("request", req),
        )

        resp, err := handler(ctx, req) // call the actual service method
        s, _ := status.FromError(err)

        logger.Info("gRPC request end",
            zap.String("method", info.FullMethod),
            zap.Duration("latency", time.Since(start)),
            zap.Int32("grpc_code", int32(s.Code())),
            zap.Any("response", resp),
            zap.Error(err),
        )

        return resp, err
    }
}
