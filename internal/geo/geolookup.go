package geo

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

type LookupService interface {
	CountryISOCode(ipStr string) (string, error)
	Close() error
}

type GeoLookupService struct {
	db *geoip2.Reader
}

func NewGeoLookupService(dbPath string) (*GeoLookupService, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}

	return &GeoLookupService{db: db}, nil
}

// CountryISOCode returns the 2-letter country code for an IP address (e.g. "US").
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

func (g *GeoLookupService) Close() error {
	return g.db.Close()
}

// Custom error for invalid IP addresses
var ErrInvalidIP = &InvalidIPError{"invalid IP address format"}

type InvalidIPError struct {
	msg string
}

func (e *InvalidIPError) Error() string {
	return e.msg
}
