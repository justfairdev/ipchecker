package server

import (
	"github.com/justfairdev/ipchecker/internal/geo"
	"github.com/justfairdev/ipchecker/internal/grpcserver"
	"github.com/justfairdev/ipchecker/internal/logger"
	"github.com/justfairdev/ipchecker/internal/middleware"
	pb "github.com/justfairdev/ipchecker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewGRPCServer constructs, configures, and returns a new gRPC server instance.
//
// This setup includes the following configurations:
//   - Structured logging using the configured Zap logger.
//   - Unary interceptor middleware for detailed logging of RPC requests and responses.
//   - Reflection service registration to support clients such as grpcurl and grpc_cli.
//   - Registration of the IPChecker service implementation for handling IP-check requests.
//
// Parameters:
//   - geoService: a GeoLookupService implementation used by the IPChecker server to perform geographic lookups.
//
// Returns:
//   - *grpc.Server: A fully configured gRPC server instance.
//   - error: An initialization error, if logger or server setup fails.
func NewGRPCServer(geoService *geo.GeoLookupService) (*grpc.Server, error) {
	// Initialize application logger with structured JSON output.
	log, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	// Create gRPC server with logging interceptor middleware for comprehensive request tracing.
	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryLoggingInterceptor(log)),
	)

	// Enable gRPC reflection to facilitate service discovery by reflection-enabled clients.
	reflection.Register(grpcSrv)

	// Instantiate and register the IPChecker service handler implementation.
	ipCheckerService := grpcserver.NewIPCheckerServer(geoService)
	pb.RegisterIPCheckerServer(grpcSrv, ipCheckerService)

	return grpcSrv, nil
}
