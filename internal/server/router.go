package server

import (
    "github.com/gin-gonic/gin"
    "github.com/justfairdev/ipchecker/internal/handler"
)

// RegisterRoutes attaches all routes to the provided gin.Engine
func RegisterRoutes(r *gin.Engine, ipChecker *handler.IPChecker) {
    // Example grouping for API version 1
    v1 := r.Group("/api/v1")

    // This is where you define your routes:
    v1.POST("/ip-check", ipChecker.CheckIP)

    // If you add more routes later, they might look like:
    // v1.POST("/some-other-endpoint", someOtherHandler)
}
