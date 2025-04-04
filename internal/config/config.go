package config

import (
	"os"
)

// Config represents the application configuration loaded from environment variables.
type Config struct {
	HTTPPort      string // Server listening port, defaults to "8080" if not specified.
	MaxMindDBPath string // Filesystem path to the MaxMind GeoLite2 database, defaults to "./GeoLite2-Country.mmdb".
}

// Load returns a Config object populated with values from environment variables.
// If an environment variable is not set, the default value is used instead.
//
// Environment Variables:
//   - HTTP_PORT: specifies the server HTTP port (default: "8080").
//   - MAXMIND_DB_PATH: specifies the file path to the MaxMind GeoLite2 database (default: "./GeoLite2-Country.mmdb").
//
// Returns:
//   - *Config: pointer to initialized Config struct.
//   - error: currently always returns nil; future revisions may include actual error handling as needed.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		MaxMindDBPath: getEnv("MAXMIND_DB_PATH", "./GeoLite2-Country.mmdb"),
	}
	return cfg, nil
}

// getEnv retrieves an environment variable using the provided key.
// If the environment variable is not set or empty, it returns the specified default value.
//
// Parameters:
//   - key (string): the environment variable key to retrieve.
//   - defaultVal (string): the default value to return if the environment variable is not found.
//
// Returns:
//   - string: the environment value if present; otherwise, the provided default value.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
