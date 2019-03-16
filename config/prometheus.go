package config

// PrometheusConfig holds configuration about Prometheus databases to backup.
type PrometheusConfig struct {
	// AdminAPIHost to contact when making snapshot.
	AdminAPIHost string `validate:"required"`

	// DataDirectory is the directory in which Prometheus data is stored.
	DataDirectory string `validate:"required"`
}
