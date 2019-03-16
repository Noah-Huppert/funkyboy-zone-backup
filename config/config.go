package config

// Config holds tool configuration.
type Config struct {
	// Upload configuration.
	Upload UploadConfig `validate:"required"`

	// Files holds configuration of standard files to backup.
	Files map[string]FileConfig

	// Prometheus holds configuration for Prometheus databases to backup.
	Prometheus map[string]PrometheusConfig
}
