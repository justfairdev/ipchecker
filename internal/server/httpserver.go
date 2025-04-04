package server

import (
	"github.com/gin-gonic/gin"
	"github.com/justfairdev/ipchecker/internal/geo"
	"github.com/justfairdev/ipchecker/internal/handler"
	"github.com/justfairdev/ipchecker/internal/logger"
	"github.com/justfairdev/ipchecker/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/justfairdev/ipchecker/docs" // Required for Swagger documentation initialization
)

// NewHTTPServer initializes and configures a new Gin HTTP server instance with custom middlewares, route handlers,
// and Swagger documentation support.
//
// The HTTP server is configured with:
//
// - Structured logging and panic recovery middleware based on the Zap logging framework.
// - A handler (`IPChecker`) for checking IP addresses against allowable country codes, leveraging the provided geographical lookup service.
// - Automated Swagger API documentation accessible at the '/swagger' endpoint for interactive exploration.
//
// Parameters:
//   - geoService: A geographical lookup implementation that the IPChecker handler utilizes for IP geolocation functionality.
//
// Returns:
//   - *gin.Engine:  Fully initialized Gin engine configured with routes, middleware, and Swagger documentation.
//   - error: Error indicating issue during logger initialization.
//
// Usage:
//
//	After running this HTTP server, you can make requests to endpoints such as:
//	POST api/v1/ip-check
//	with body:
//	{
//	    "ip_address": "128.101.101.101",
//	    "allowed_countries": ["US", "CA"]
//	}
//
// Interactive Swagger API documentation is available at:
//
//	http://localhost:8080/swagger/index.html
//
// Similarly, you can test gRPC service using `grpcurl`:
//
//	grpcurl -plaintext -d '{"ip_address":"128.101.101.101","allowed_countries":["US","CA"]}' \
//	  localhost:50051 ipchecker.v1.IPChecker/CheckIP
func NewHTTPServer(geoService *geo.GeoLookupService) (*gin.Engine, error) {
	// Initialize structured Zap logger for consistent and reliable request tracing
	log, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	// Instantiate Gin router without default middlewares for more control
	r := gin.New()

	// Attach customized middleware for structured logging and panic recovery
	r.Use(
		middleware.GinLogger(log),
		middleware.GinRecovery(log),
	)

	// Initialize the IPChecker route handler with the geo lookup service dependency
	ipChecker := handler.NewIPChecker(geoService)

	// Register IPChecker routes to the Gin server
	RegisterRoutes(r, ipChecker)

	// Optionally enable Swagger UI at '/swagger' for convenient API testing and documentation viewing
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r, nil
}
