package server  

import (  
    "github.com/gin-gonic/gin"  
    "github.com/justfairdev/ipchecker/internal/handler"  
)  

// RegisterRoutes sets up and attaches HTTP endpoints (routes) to the provided Gin engine.  
//  
// This function organizes endpoints within logically grouped route prefixes,  
// such as API versioning paths for maintainability and clarity.  
//  
// Parameters:  
//   - r: The Gin HTTP engine instance to which the routes will be attached.  
//   - ipChecker: An instance of the IPChecker handler responsible for handling IP-checking requests.  
//  
// Current endpoints registered:  
//   - POST /api/v1/ip-check : Verifies whether an IP address is within a list of allowed country codes.  
//  
// Example JSON request payload:  
//   {  
//       "ip_address": "128.101.101.101",  
//       "allowed_countries": ["US", "CA"]  
//   }  
//  
// Future endpoints can be efficiently added within this function following the existing structure,  
// ensuring ease of management and readability.  
func RegisterRoutes(r *gin.Engine, ipChecker *handler.IPChecker) {  
    // Group routes under API Version 1 prefix for version control and structured endpoint management.  
    v1 := r.Group("/api/v1")  

    // IP address checking route.  
    v1.POST("/ip-check", ipChecker.CheckIP)  

    // Additional API routes may be defined here as needed.  
    // Example:  
    // v1.POST("/another-endpoint", anotherHandler.Method)  
}  