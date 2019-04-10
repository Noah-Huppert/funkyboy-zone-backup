package config

// MetricsConfig hold configuration for metrics pushed to Prometheus.
type MetricsConfig struct {
	// Enabled indicates if metrics should be pushed.
	Enabled bool

	// PushGatewayHost is the host at which the Prometheus Push Gateway that metrics will be pushed to can be
	// access. Must include a scheme.
	PushGatewayHost string `validate:"url,required" default:"http://localhost:9091"`

	// LabelHost is the value of the `host` label in pushed metrics
	LabelHost string `validate:"required" default:"mountain-backup"`
}
