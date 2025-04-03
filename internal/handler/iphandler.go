package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justfairdev/ipchecker/internal/geo"
)

type IPCheckRequest struct {
	IPAddress        string   `json:"ip_address" binding:"required"`
	AllowedCountries []string `json:"allowed_countries" binding:"required"`
}

type IPCheckResponse struct {
	Allowed bool   `json:"allowed"`
	Country string `json:"country"`
}

type IPChecker struct {
	geoService *geo.GeoLookupService
}

func NewIPChecker(geoService *geo.GeoLookupService) *IPChecker {
	return &IPChecker{geoService: geoService}
}

func (c *IPChecker) CheckIP(ctx *gin.Context) {
	var req IPCheckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	countryCode, err := c.geoService.CountryISOCode(req.IPAddress)
	if err != nil {
		if err == geo.ErrInvalidIP {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid IP address"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unable to lookup country"})
		}
		return
	}

	// Check if the country's ISO code is in the AllowedCountries list.
	allowed := false
	for _, ac := range req.AllowedCountries {
		if ac == countryCode {
			allowed = true
			break
		}
	}

	ctx.JSON(http.StatusOK, IPCheckResponse{
		Allowed: allowed,
		Country: countryCode,
	})
}
