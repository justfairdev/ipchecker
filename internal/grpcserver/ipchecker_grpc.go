package grpcserver

import (
	"context"

	"github.com/justfairdev/ipchecker/internal/geo"
	pb "github.com/justfairdev/ipchecker/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IPCheckerServerImpl implements the IPChecker gRPC service defined in the protobuf specification.
// This server handles IP address checks against geographical locations based on provided criteria.
type IPCheckerServerImpl struct {
	pb.UnimplementedIPCheckerServer
	geoService geo.LookupService
}

// NewIPCheckerServer constructs a new IPCheckerServerImpl instance with the provided geographical lookup service.
//
// Parameters:
//   - gs: An implementation of geo.LookupService for geographical IP address resolution.
//
// Returns:
//   - Pointer to IPCheckerServerImpl configured with the specified geo service.
func NewIPCheckerServer(gs geo.LookupService) *IPCheckerServerImpl {
	return &IPCheckerServerImpl{geoService: gs}
}

// CheckIP processes the IPCheckRequest by performing a geographical lookup of the specified IP address
// and verifies if it originates from one of the allowed countries provided.
//
// Parameters:
//   - ctx: Context carrying metadata and deadlines for the request handling lifecycle.
//   - req: IPCheckRequest containing the target IP address and list of allowed ISO 3166-1 alpha-2 country codes.
//
// Returns:
//   - *pb.IPCheckResponse: Contains the country code associated with the IP and whether it is permitted.
//   - error: Returns a gRPC status error if IP lookup fails or the IP address format is invalid.
func (s *IPCheckerServerImpl) CheckIP(ctx context.Context, req *pb.IPCheckRequest) (*pb.IPCheckResponse, error) {
	// Perform geographical lookup to obtain the country associated with the provided IP address.
	country, err := s.geoService.CountryISOCode(req.GetIpAddress())
	if err != nil {
		// Return a gRPC error indicating the provided IP address is invalid or geo lookup failed.
		return nil, status.Errorf(codes.InvalidArgument, "invalid IP address: %v", err)
	}

	// Check whether the resolved country code is within the allowed countries list.
	allowed := false
	for _, allowedCountry := range req.GetAllowedCountries() {
		if allowedCountry == country {
			allowed = true
			break
		}
	}

	// Return the result indicating if the IP is allowed and its associated country code.
	return &pb.IPCheckResponse{
		Allowed: allowed,
		Country: country,
	}, nil
}
