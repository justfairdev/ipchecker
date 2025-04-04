package geo


// MockGeoLookupService is a mock implementation of the LookupService interface.
type MockGeoLookupService struct {
    MockCountryCode string
    MockError       error
}

// NewMockGeoLookupService creates a new mock that returns the given country code and error.
func NewMockGeoLookupService(countryCode string, err error) *MockGeoLookupService {
    return &MockGeoLookupService{
        MockCountryCode: countryCode,
        MockError:       err,
    }
}

// CountryISOCode returns the configured country code or error.
func (m *MockGeoLookupService) CountryISOCode(ipAddress string) (string, error) {
    if m.MockError != nil {
        return "", m.MockError
    }
    return m.MockCountryCode, nil
}

// Close is a no-op for the mock.
func (m *MockGeoLookupService) Close() error {
    return nil
}
