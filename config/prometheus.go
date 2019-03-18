package config

// PrometheusConfig holds configuration about Prometheus databases to backup.
type PrometheusConfig struct {
	// AdminAPIHost to contact when making snapshot. Must include a scheme.
	AdminAPIHost string `validate:"required" validate:"url"`

	// DataDirectory is the directory in which Prometheus data is stored.
	DataDirectory string `validate:"required"`
}
