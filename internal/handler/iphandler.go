package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/justfairdev/ipchecker/internal/dtos"
    "github.com/justfairdev/ipchecker/internal/geo"
)

// IPChecker holds the dependencies needed by this handler.
type IPChecker struct {
    geoService geo.LookupService
}

// NewIPChecker is a constructor that injects the geo service (using the interface).
func NewIPChecker(geoService geo.LookupService) *IPChecker {
    return &IPChecker{geoService: geoService}
}

// CheckIP godoc
// @Summary      Check if IP is in allowed countries
// @Description  Takes an IP address & list of allowed countries, returns whether it's allowed.
// @Tags         IP
// @Accept       json
// @Produce      json
// @Param        requestBody body dtos.IPCheckRequest true "IP Check Payload"
// @Success      200  {object}  dtos.IPCheckResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /ip-check [post]
func (c *IPChecker) CheckIP(ctx *gin.Context) {
    var req dtos.IPCheckRequest
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

    allowed := false
    for _, ac := range req.AllowedCountries {
        if ac == countryCode {
            allowed = true
            break
        }
    }

    ctx.JSON(http.StatusOK, dtos.IPCheckResponse{
        Allowed: allowed,
        Country: countryCode,
    })
}
