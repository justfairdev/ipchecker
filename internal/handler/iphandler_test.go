package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/justfairdev/ipchecker/internal/geo"
	"github.com/justfairdev/ipchecker/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestIPChecker_CheckIP_Success(t *testing.T) {
	// Set Gin to test mode.
	gin.SetMode(gin.TestMode)

	// Create a mock geo service that always returns "US" with no error.
	mockGeo := geo.NewMockGeoLookupService("US", nil)

	// Create a new instance of the IPChecker handler using the mock.
	ipChecker := handler.NewIPChecker(mockGeo)

	// Create a Gin engine and register the route.
	router := gin.Default()
	router.POST("/ip-check", ipChecker.CheckIP)

	// Prepare a valid JSON request.
	reqBody := `{"ip_address":"128.101.101.101","allowed_countries":["US","CA"]}`
	req, err := http.NewRequest(http.MethodPost, "/ip-check", strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the response.
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert that the response status is 200 OK.
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse and assert the response body.
	// Assuming dtos.IPCheckResponse is similar to:
	// type IPCheckResponse struct {
	//     Allowed bool   `json:"allowed"`
	//     Country string `json:"country"`
	// }
	responseBody := recorder.Body.String()
	assert.Contains(t, responseBody, `"allowed":true`)
	assert.Contains(t, responseBody, `"country":"US"`)
}

func TestIPChecker_CheckIP_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockGeo := geo.NewMockGeoLookupService("US", nil)
	ipChecker := handler.NewIPChecker(mockGeo)
	router := gin.Default()
	router.POST("/ip-check", ipChecker.CheckIP)

	// Provide invalid JSON.
	reqBody := `{"ip_address": "128.101.101.101", "allowed_countries": }`
	req, err := http.NewRequest(http.MethodPost, "/ip-check", strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Expect 400 Bad Request.
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
