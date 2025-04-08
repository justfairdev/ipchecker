package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justfairdev/ipchecker/internal/dtos"
	"github.com/justfairdev/ipchecker/internal/geo"
)

// IPChecker provides HTTP handlers for IP address verification against allowed countries.
type IPChecker struct {
	geoService geo.LookupService
}

// NewIPChecker constructs a new IPChecker handler with the given Geo lookup service dependency.
//
// Parameters:
//   - geoService: An implementation of geo.LookupService used to determine the country of IP addresses.
//
// Returns:
//   - *IPChecker: A pointer to the initialized IPChecker handler instance.
func NewIPChecker(geoService geo.LookupService) *IPChecker {
	return &IPChecker{geoService: geoService}
}

// CheckIP godoc
// @Summary      Verify if an IP address originates from allowed countries.
// @Description  Accepts an IP address and a list of allowed countries; returns whether the IP address is permitted based on its location.
// @Tags         IP
// @Accept       json
// @Produce      json
// @Param        requestBody body dtos.IPCheckRequest true "IP check request payload."
// @Success      200 {object} dtos.IPCheckResponse "Successful IP check operation."
// @Failure      400 {object} map[string]string "Invalid request payload or malformed IP address."
// @Failure      500 {object} map[string]string "Internal server error during IP geolocation lookup."
// @Router       /ip-check [post]
func (c *IPChecker) CheckIP(ctx *gin.Context) {
	var req dtos.IPCheckRequest

	// Bind the incoming JSON request payload to the IPCheckRequest struct.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Perform IP classification using the provided geo lookup service.
	countryCode, err := c.geoService.CountryISOCode(req.IPAddress)
	if err != nil {
		// Handle specific invalid IP format error explicitly.
		if err == geo.ErrInvalidIP {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid IP address"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unable to lookup country"})
		}
		return
	}

	// Determine if the resolved IP country code is within the allowed countries.
	allowed := false
	for _, allowedCountry := range req.AllowedCountries {
		if allowedCountry == countryCode {
			allowed = true
			break
		}
	}

	// Return a structured JSON response indicating the IP address permission status and country code.
	ctx.JSON(http.StatusOK, dtos.IPCheckResponse{
		Allowed: allowed,
		Country: countryCode,
	})
}
