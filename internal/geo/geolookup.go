package geo

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

// LookupService defines methods for IP-based geolocation queries.
type LookupService interface {
	// CountryISOCode retrieves the ISO 3166-1 alpha-2 country code (e.g., "US") for the given IP address.
	CountryISOCode(ipStr string) (string, error)

	// Close safely releases any underlying resources associated with the LookupService.
	Close() error
}

// GeoLookupService implements the LookupService interface using the MaxMind GeoIP2 database.
type GeoLookupService struct {
	db *geoip2.Reader
}

// NewGeoLookupService initializes and returns a new GeoLookupService instance.
// It opens the GeoIP2 database located at the specified file path.
//
// Parameters:
//   - dbPath: File system path to the MaxMind GeoLite2 or GeoIP2 database file.
//
// Returns:
//   - *GeoLookupService: An initialized GeoLookupService instance.
//   - error: Error if the database could not be opened successfully.
func NewGeoLookupService(dbPath string) (*GeoLookupService, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}

	return &GeoLookupService{db: db}, nil
}

// CountryISOCode takes an IP address string and returns the corresponding two-letter country ISO code.
// It is suitable for quick lookup operations.
//
// Parameters:
//   - ipStr: String representation of the IP address to be checked.
//
// Returns:
//   - string: The ISO 3166-1 alpha-2 country code (e.g., "US") associated with the provided IP.
//   - error: An error if the IP format is invalid, or if the lookup operation fails.
func (g *GeoLookupService) CountryISOCode(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", ErrInvalidIP
	}

	record, err := g.db.Country(ip)
	if err != nil {
		return "", err
	}

	return record.Country.IsoCode, nil
}

// Close releases the internal resources used by the GeoLookupService.
// This should be called when the service is no longer needed to avoid resource leaks.
//
// Returns:
//   - error: An error if closing the database resource fails.
func (g *GeoLookupService) Close() error {
	return g.db.Close()
}

// ErrInvalidIP represents an error returned when the provided IP address is incorrectly formatted.
var ErrInvalidIP = &InvalidIPError{"invalid IP address format"}

// InvalidIPError indicates an error encountered during IP parsing due to invalid format.
type InvalidIPError struct {
	msg string
}

// Error satisfies the error interface for InvalidIPError and returns the error message.
//
// Returns:
//   - string: Descriptive error message indicating why the IP was considered invalid.
func (e *InvalidIPError) Error() string {
	return e.msg
}
