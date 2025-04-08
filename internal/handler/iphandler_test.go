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

// TestIPChecker_CheckIP_Success ensures that the IPChecker handler responds correctly
// when provided valid input data and the geographical lookup service returns a successful result.
func TestIPChecker_CheckIP_Success(t *testing.T) {
	// Set Gin's running mode to TestMode for predictable testing behavior.
	gin.SetMode(gin.TestMode)

	// Initialize a mock GeoLookupService to simulate successful geolocation lookup,
	// always returning "US" as the country code.
	mockGeo := geo.NewMockGeoLookupService("US", nil)

	// Instantiate the IPChecker handler using the mocked GeoLookupService.
	ipChecker := handler.NewIPChecker(mockGeo)

	// Configure Gin router with the IP check handler route.
	router := gin.Default()
	router.POST("/ip-check", ipChecker.CheckIP)

	// Prepare test HTTP POST request body with valid input JSON.
	reqBody := `{"ip_address":"128.101.101.101","allowed_countries":["US","CA"]}`
	req, err := http.NewRequest(http.MethodPost, "/ip-check", strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request and record the response for validation.
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Validate that the response status code is HTTP 200 OK.
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Verify that the returned JSON response contains correct values indicating successful lookup.
	responseBody := recorder.Body.String()
	assert.Contains(t, responseBody, `"allowed":true`)
	assert.Contains(t, responseBody, `"country":"US"`)
}

// TestIPChecker_CheckIP_InvalidJSON ensures the IPChecker handler returns an HTTP 400 Bad Request status
// when it receives improperly formatted JSON input.
func TestIPChecker_CheckIP_InvalidJSON(t *testing.T) {
	// Set Gin's running mode to TestMode.
	gin.SetMode(gin.TestMode)

	// Initialize the mock GeoLookupService; its response is irrelevant for invalid JSON inputs.
	mockGeo := geo.NewMockGeoLookupService("US", nil)

	// Instantiate the IPChecker handler with the mocked GeoLookupService.
	ipChecker := handler.NewIPChecker(mockGeo)

	// Configure Gin router for handling IP checker requests.
	router := gin.Default()
	router.POST("/ip-check", ipChecker.CheckIP)

	// Prepare test HTTP POST request body containing invalid (malformed) JSON.
	reqBody := `{"ip_address": "128.101.101.101", "allowed_countries": }`
	req, err := http.NewRequest(http.MethodPost, "/ip-check", strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request and record the handler response.
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Validate that an HTTP 400 Bad Request status code is returned due to invalid JSON syntax.
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
