package server

import (
    "github.com/justfairdev/ipchecker/internal/geo"
    "github.com/justfairdev/ipchecker/internal/grpcserver"
	"github.com/justfairdev/ipchecker/internal/middleware"
    "github.com/justfairdev/ipchecker/internal/logger"
    pb "github.com/justfairdev/ipchecker/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

// NewGRPCServer creates a gRPC server, registers reflection & the IP checker service.
func NewGRPCServer(geoService *geo.GeoLookupService) (*grpc.Server, error) {
    // 1) Create logger
    log, err := logger.NewLogger()
    if err != nil {
        return nil, err
    }

    // 2) Construct the gRPC server (optionally attach interceptors)
    grpcSrv := grpc.NewServer(
        grpc.UnaryInterceptor(middleware.UnaryLoggingInterceptor(log)),
    )

    // 3) Optionally enable reflection so grpcurl & others can discover services
    reflection.Register(grpcSrv)

    // 4) Register your actual service implementation
    impl := grpcserver.NewIPCheckerServer(geoService)
    pb.RegisterIPCheckerServer(grpcSrv, impl)

    return grpcSrv, nil
}
