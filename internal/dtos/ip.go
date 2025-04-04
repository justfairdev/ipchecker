package dtos

// IPCheckRequest represents the request payload for checking if a given IP address
// originates from one of the allowed countries.
//
// swagger:model IPCheckRequest
type IPCheckRequest struct {
	// IPAddress is the IP address to be verified.
	// Required field.
	IPAddress string `json:"ip_address" binding:"required"`

	// AllowedCountries is a list of acceptable country codes (ISO 3166-1 alpha-2 format).
	// The IP address must originate from one of these countries.
	// Required field.
	AllowedCountries []string `json:"allowed_countries" binding:"required"`
}

// IPCheckResponse represents the response payload after checking the requested IP address.
//
// swagger:model IPCheckResponse
type IPCheckResponse struct {
	// Allowed indicates whether the given IP address is from one of the allowed countries.
	Allowed bool `json:"allowed"`

	// Country is the ISO 3166-1 alpha-2 country code associated with the IP address.
	Country string `json:"country"`
}
