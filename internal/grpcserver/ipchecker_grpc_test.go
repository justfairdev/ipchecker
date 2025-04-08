package grpcserver_test

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/justfairdev/ipchecker/internal/geo"
	"github.com/justfairdev/ipchecker/internal/grpcserver"
	pb "github.com/justfairdev/ipchecker/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// dialer returns a connection dialing function used for creating in-memory gRPC connections during testing.
//
// Parameters:
//   - listener: a reference to a bufconn listener.
//
// Returns:
//   - func(context.Context, string) (net.Conn, error): a dialer function suitable for grpc.DialContext.
func dialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, address string) (net.Conn, error) {
		return listener.Dial()
	}
}

// TestIPCheckerGRPC_CheckIP_Success verifies successful behavior of the gRPC IPChecker Server
// when the geo lookup service returns a valid country code.
func TestIPCheckerGRPC_CheckIP_Success(t *testing.T) {
	// Initialize an in-memory gRPC server using the bufconn package.
	listener := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()

	// Create a mock geographical lookup service configured to return "US" as the country code.
	mockGeo := geo.NewMockGeoLookupService("US", nil)

	// Instantiate the IPChecker gRPC server implementation with the mock service.
	ipCheckerSvc := grpcserver.NewIPCheckerServer(mockGeo)
	pb.RegisterIPCheckerServer(grpcServer, ipCheckerSvc)

	// Serve the gRPC server concurrently for the duration of this test.
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to start gRPC test server: %v", err)
		}
	}()

	// Establish a client connection to the in-memory gRPC server.
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer(listener)), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewIPCheckerClient(conn)

	// Prepare a well-formed test request.
	req := &pb.IPCheckRequest{
		IpAddress:        "128.101.101.101",
		AllowedCountries: []string{"US", "CA"},
	}

	// Execute the gRPC method under testing.
	resp, err := client.CheckIP(ctx, req)

	// Verify expectations and response correctness.
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Allowed, "Expected the IP address to be allowed.")
	assert.Equal(t, "US", resp.Country, "Expected the IP country code to be 'US'.")
}

// TestIPCheckerGRPC_CheckIP_Failure verifies proper handling of an error scenario
// when the geo lookup service encounters an internal error.
func TestIPCheckerGRPC_CheckIP_Failure(t *testing.T) {
	// Initialize an in-memory gRPC test server using bufconn.
	listener := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()

	// Configure the mock geographical lookup service to always return an error.
	mockGeo := geo.NewMockGeoLookupService("", errors.New("geo service error"))

	// Instantiate the IPChecker gRPC server implementation with the failing mock service.
	ipCheckerSvc := grpcserver.NewIPCheckerServer(mockGeo)
	pb.RegisterIPCheckerServer(grpcServer, ipCheckerSvc)

	// Run the gRPC server concurrently for the test.
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("gRPC test server exited with error: %v", err)
		}
	}()

	// Connect to the in-memory gRPC server as a client.
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer(listener)), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewIPCheckerClient(conn)

	// Send a valid request and expect an error response due to mockGeo's configuration.
	req := &pb.IPCheckRequest{
		IpAddress:        "128.101.101.101",
		AllowedCountries: []string{"US", "CA"},
	}

	// Execute the CheckIP call and assert error conditions.
	resp, err := client.CheckIP(ctx, req)
	assert.Error(t, err, "Expected error due to simulated geo service error.")
	assert.Nil(t, resp, "Expected no response due to internal geo service failure.")
}
