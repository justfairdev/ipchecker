package geo

// MockGeoLookupService provides a mock implementation of the LookupService interface,
// primarily intended for use in unit tests.
type MockGeoLookupService struct {
	// MockCountryCode is the predefined ISO 3166-1 alpha-2 country code to return during testing.
	MockCountryCode string

	// MockError is the predefined error to return when simulating error scenarios.
	MockError error
}

// NewMockGeoLookupService initializes a new MockGeoLookupService with the specified
// mock country code and error.
//
// Parameters:
//   - countryCode: The mock ISO 3166-1 alpha-2 country code to be returned (e.g., "US").
//   - err: The mock error to return for simulating error conditions. Set to nil for normal operation.
//
// Returns:
//   - *MockGeoLookupService: An initialized mock service instance.
func NewMockGeoLookupService(countryCode string, err error) *MockGeoLookupService {
	return &MockGeoLookupService{
		MockCountryCode: countryCode,
		MockError:       err,
	}
}

// CountryISOCode simulates retrieving a country ISO code for a given IP address.
// Returns either the configured mock country code or the configured mock error.
//
// Parameters:
//   - ipAddress: The IP address string to check (ignored by the mock implementation).
//
// Returns:
//   - string: The predefined mock country code, if no mock error is specified.
//   - error: The predefined mock error, if any; otherwise nil.
func (m *MockGeoLookupService) CountryISOCode(ipAddress string) (string, error) {
	if m.MockError != nil {
		return "", m.MockError
	}
	return m.MockCountryCode, nil
}

// Close is a mock implementation to satisfy the LookupService interface.
// It performs no operation and always returns nil.
//
// Returns:
//   - error: Always nil, representing successful resource release.
func (m *MockGeoLookupService) Close() error {
	return nil
}
