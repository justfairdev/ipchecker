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

// AppServer holds both the HTTP server (Gin) and the gRPC server (along with the geoService).
type AppServer struct {
	HTTPServer *gin.Engine
	GRPCServer *grpc.Server
	geoService *geo.GeoLookupService
}

// NewAppServer creates an AppServer that can serve both HTTP (Gin) and gRPC.
func NewAppServer(cfg *config.Config) (*AppServer, error) {
	// 1) Create the Geo service once
	geoSvc, err := geo.NewGeoLookupService(cfg.MaxMindDBPath)
	if err != nil {
		return nil, err
	}

	// 2) Build HTTP server
	httpServer, err := NewHTTPServer(geoSvc)
	if err != nil {
		return nil, err
	}

	// 3) Build gRPC server
	grpcSrv, err := NewGRPCServer(geoSvc)
	if err != nil {
		return nil, err
	}

	return &AppServer{
		HTTPServer: httpServer,
		GRPCServer: grpcSrv,
		geoService: geoSvc,
	}, nil
}

// Start runs the HTTP server on httpPort AND the gRPC server on grpcPort (concurrently).
func (s *AppServer) Start(httpPort, grpcPort string) error {
	// Start gRPC in a goroutine
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
		if err != nil {
			log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
		}
		log.Printf("gRPC server listening on :%s", grpcPort)

		if err := s.GRPCServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Run the HTTP server (blocking)
	log.Printf("HTTP server listening on :%s", httpPort)
	return s.HTTPServer.Run(":" + httpPort)
}

// Stop gracefully stops gRPC and closes the geo DB (you can also shut down HTTP if desired).
func (s *AppServer) Stop() {
	log.Println("Stopping gRPC server gracefully...")
	s.GRPCServer.GracefulStop()

	log.Println("Closing GeoLookupService...")
	s.geoService.Close()
}
