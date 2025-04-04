package grpcserver_test

import (
	
	"context"
	"errors"
	"net"
	"testing"
	"log"

	"github.com/justfairdev/ipchecker/internal/geo"
	"github.com/justfairdev/ipchecker/internal/grpcserver"
	pb "github.com/justfairdev/ipchecker/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// dialer returns a function that dials a bufconn listener.
func dialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, address string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestIPCheckerGRPC_CheckIP_Success(t *testing.T) {
	// Create an in-memory gRPC server using bufconn.
	listener := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()

	// Create a mock geo service that returns "US" (simulate success).
	mockGeo := geo.NewMockGeoLookupService("US", nil)
	ipCheckerSvc := grpcserver.NewIPCheckerServer(mockGeo)
	pb.RegisterIPCheckerServer(grpcServer, ipCheckerSvc)

	// Start the server in a goroutine.
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to start servers: %v", err)
		}
	}()

	// Create a client connection to the in-memory server.
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer(listener)), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewIPCheckerClient(conn)

	// Create a valid gRPC request.
	req := &pb.IPCheckRequest{
		IpAddress:        "128.101.101.101",
		AllowedCountries: []string{"US", "CA"},
	}
	resp, err := client.CheckIP(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Allowed)
	assert.Equal(t, "US", resp.Country)
}

func TestIPCheckerGRPC_CheckIP_Failure(t *testing.T) {
	// Create an in-memory gRPC server using bufconn.
	listener := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()

	// Create a mock geo service that always returns an error.
	mockGeo := geo.NewMockGeoLookupService("", errors.New("geo service error"))
	ipCheckerSvc := grpcserver.NewIPCheckerServer(mockGeo)
	pb.RegisterIPCheckerServer(grpcServer, ipCheckerSvc)

	// Start the server in a goroutine.
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// Create a client connection to the in-memory server.
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer(listener)), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewIPCheckerClient(conn)

	// Create a request; expect an error due to the mock.
	req := &pb.IPCheckRequest{
		IpAddress:        "128.101.101.101",
		AllowedCountries: []string{"US", "CA"},
	}
	resp, err := client.CheckIP(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
