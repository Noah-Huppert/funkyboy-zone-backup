package config

// Config holds tool configuration
type Config struct {
	// DigtalOcean configuration
	DigitalOcean DigitalOceanConfig `validate:"required"`
}

// DigitalOceanConfig holds Digital Ocean configuration
type DigitalOceanConfig struct {
	// KeyID for spaces API
	KeyID string `validate:"required"`

	// AccessKey for spaces API
	AccessKey string `validate:"required"`
}
