package config

import (
	"os"
)

type Config struct {
	HTTPPort      string
	MaxMindDBPath string
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		MaxMindDBPath: getEnv("MAXMIND_DB_PATH", "./GeoLite2-Country.mmdb"),
	}
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
