package dtos

// IPCheckRequest represents the incoming JSON payload for an IP check.
//
// swagger:model IPCheckRequest
type IPCheckRequest struct {
	IPAddress        string   `json:"ip_address" binding:"required"`
	AllowedCountries []string `json:"allowed_countries" binding:"required"`
}

// IPCheckResponse represents the outgoing JSON payload for an IP check.
//
// swagger:model IPCheckResponse
type IPCheckResponse struct {
	Allowed bool   `json:"allowed"`
	Country string `json:"country"`
}
