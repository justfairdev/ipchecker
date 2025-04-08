package server

import (
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/justfairdev/ipchecker/internal/config"
	"github.com/justfairdev/ipchecker/internal/geo"
	"google.golang.org/grpc"
)

// AppServer encapsulates both HTTP and gRPC server instances along with the shared GeoLookupService dependency.
//
// This struct allows simultaneous management of the HTTP REST endpoints (via Gin) and gRPC endpoints, providing a unified
// approach to serving multiple types of clients with shared underlying services.
type AppServer struct {
	HTTPServer *gin.Engine           // Instance of the Gin-powered HTTP server
	GRPCServer *grpc.Server          // Instance of the gRPC server
	geoService *geo.GeoLookupService // Shared GeoLookup service instance used by both servers
}

// NewAppServer initializes an AppServer instance configured for both HTTP and gRPC servers.
//
// The initialization process involves:
//   - Creating a single shared GeoLookupService instance with the specified MaxMind database.
//   - Constructing and configuring the Gin HTTP server with routes, middleware, and handlers.
//   - Constructing and configuring the gRPC server instance with appropriate service handlers.
//
// Parameters:
//   - cfg: A configuration struct containing critical parameters (e.g., path to MaxMind Geo database).
//
// Returns:
//   - *AppServer: A fully initialized AppServer instance ready for operation.
//   - error: If initialization fails, returns an error describing the issue.
func NewAppServer(cfg *config.Config) (*AppServer, error) {
	// Initialize shared GeoLookupService dependency
	geoSvc, err := geo.NewGeoLookupService(cfg.MaxMindDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GeoLookupService: %w", err)
	}

	// Initialize and configure HTTP server (Gin engine)
	httpServer, err := NewHTTPServer(geoSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize HTTP server: %w", err)
	}

	// Initialize and configure gRPC server
	grpcSrv, err := NewGRPCServer(geoSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gRPC server: %w", err)
	}

	// Return the fully configured AppServer instance
	return &AppServer{
		HTTPServer: httpServer,
		GRPCServer: grpcSrv,
		geoService: geoSvc,
	}, nil
}

// Start concurrently launches the HTTP server and the gRPC server, handling requests on their respective ports.
//
// Execution flow:
//   - gRPC server startup occurs asynchronously in a separate goroutine.
//   - HTTP server startup occurs on the main thread and blocks until stopped.
//
// Parameters:
//   - httpPort: TCP port number for HTTP server (e.g., "8080").
//   - grpcPort: TCP port number for gRPC server (e.g., "50051").
//
// Returns:
//   - error: If the HTTP server fails during runtime initialization, it returns an error.
//     (gRPC server initialization errors will result in process exit via log.Fatal within the goroutine.)
func (s *AppServer) Start(httpPort, grpcPort string) error {
	// Start the gRPC server in its own goroutine concurrently with HTTP server
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port %s: %v", grpcPort, err)
		}
		log.Printf("gRPC server is running and listening on port :%s", grpcPort)

		if serveErr := s.GRPCServer.Serve(listener); serveErr != nil {
			log.Fatalf("Failed to serve the gRPC server: %v", serveErr)
		}
	}()

	// Start the HTTP server; this call is blocking
	log.Printf("HTTP server is running and listening on port :%s", httpPort)
	return s.HTTPServer.Run(":" + httpPort)
}

// Stop performs a graceful shutdown of the gRPC server and closes related services.
//
// This method ensures:
//   - Graceful stopping of the gRPC server, allowing ongoing operations to complete.
//   - Proper closure of the GeoLookupService handle (releasing database resources).
//
// Note that the HTTP server (Gin engine) currently does not have explicit graceful shutdown logic in this method.
// Developers may choose to add HTTP server graceful shutdown support if needed in the future.
func (s *AppServer) Stop() {
	log.Println("Initiating graceful shutdown of gRPC server...")
	s.GRPCServer.GracefulStop()

	log.Println("Closing GeoLookupService database connection...")
	s.geoService.Close()
}
