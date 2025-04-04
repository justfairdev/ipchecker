package grpcserver

import (
	"context"

	"github.com/justfairdev/ipchecker/internal/geo"
	pb "github.com/justfairdev/ipchecker/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IPCheckerServerImpl implements the gRPC IPChecker service.
type IPCheckerServerImpl struct {
	pb.UnimplementedIPCheckerServer
	geoService geo.LookupService
}

// NewIPCheckerServer returns a pointer to IPCheckerServerImpl.
// It accepts any type that satisfies the geo.LookupService interface.
func NewIPCheckerServer(gs geo.LookupService) *IPCheckerServerImpl {
	return &IPCheckerServerImpl{geoService: gs}
}

// CheckIP is the actual method that runs your business logic.
func (s *IPCheckerServerImpl) CheckIP(ctx context.Context, req *pb.IPCheckRequest) (*pb.IPCheckResponse, error) {
	// Perform a geo lookup using the provided interface.
	country, err := s.geoService.CountryISOCode(req.GetIpAddress())
	if err != nil {
		// Return an appropriate gRPC error code if lookup fails.
		return nil, status.Errorf(codes.InvalidArgument, "invalid IP address")
	}

	// Determine if the country's ISO code is in the allowed list.
	allowed := false
	for _, ac := range req.GetAllowedCountries() {
		if ac == country {
			allowed = true
			break
		}
	}

	// Return the response.
	return &pb.IPCheckResponse{
		Allowed: allowed,
		Country: country,
	}, nil
}
