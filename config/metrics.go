package config

// MetricsConfig hold configuration for metrics pushed to Prometheus.
type MetricsConfig struct {
	// PushGatewayHost is the host at which the Prometheus Push Gateway that metrics will be pushed to can be
	// access. Must include a scheme.
	PushGatewayHost string `validate:"url,required"`

	// LabelHost is the value of the `host` label in pushed metrics
	LabelHost string `validate:"required"`
}
