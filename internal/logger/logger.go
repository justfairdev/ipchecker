package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates and configures a new application logger using the Uber Zap library.
//
// By default, the logger is configured to:
//   - Output structured logs in JSON format.
//   - Set the minimum log severity level to "Info".
//   - Write standard logs to standard output (stdout).
//   - Write error logs to standard error (stderr).
//
// Consideration:
//
//	For production readiness, the logger configuration options—such as log format, log level, and output destinations—
//	should be externalized, typically through environment variables or a configuration file.
//
// Returns:
//   - *zap.Logger: A configured Zap logger ready for use throughout the application.
//   - error: An error if logger initialization fails due to misconfiguration.
func NewLogger() (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding:         "json",                                  // Structured logging in JSON format.
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel), // Set minimum logging level to Info.
		OutputPaths:      []string{"stdout"},                      // Standard output stream for normal log entries.
		ErrorOutputPaths: []string{"stderr"},                      // Error logs go to standard error output.
		EncoderConfig:    zap.NewProductionEncoderConfig(),        // Use recommended production encoder settings.
	}

	return cfg.Build()
}
