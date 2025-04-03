package server

import (
	"github.com/gin-gonic/gin"
	"github.com/justfairdev/ipchecker/internal/config"
	"github.com/justfairdev/ipchecker/internal/geo"
	"github.com/justfairdev/ipchecker/internal/handler"
)

func NewServer(cfg *config.Config) (*gin.Engine, error) {
	// Create a GeoLookupService using the path in config
	geoService, err := geo.NewGeoLookupService(cfg.MaxMindDBPath)
	if err != nil {
		return nil, err
	}
	// In production, remember to close this DB at shutdown.

	// Create the IP checker handler, passing in the geoService
	ipChecker := handler.NewIPChecker(geoService)

	// Create a Gin router with default middleware (logging, recovery)
	r := gin.Default()

	// 4) Register routes in a separate function
	RegisterRoutes(r, ipChecker)

	return r, nil
}
