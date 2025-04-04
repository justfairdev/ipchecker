package server

import (
	"github.com/gin-gonic/gin"
	"github.com/justfairdev/ipchecker/internal/geo"
	"github.com/justfairdev/ipchecker/internal/handler"
	"github.com/justfairdev/ipchecker/internal/logger"
	"github.com/justfairdev/ipchecker/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/justfairdev/ipchecker/docs" // Needed for Swagger docs
)

// NewHTTPServer initializes a Gin engine with routes, middlewares, and Swagger.
// geoService is passed so route handlers can share it (for IP lookups, etc.).
func NewHTTPServer(geoService *geo.GeoLookupService) (*gin.Engine, error) {
	// 1) Create a Zap logger (or re-use an existing if you prefer)
	log, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	// 2) Create a new Gin router without default middlewares
	r := gin.New()

	// 3) Attach your custom middlewares (logging, recovery, etc.)
	r.Use(
		middleware.GinLogger(log),
		middleware.GinRecovery(log),
	)

	// 4) Create handler(s)
	ipChecker := handler.NewIPChecker(geoService)

	// 5) Register routes
	RegisterRoutes(r, ipChecker)

	// 6) Optionally serve Swagger at /swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r, nil
}
