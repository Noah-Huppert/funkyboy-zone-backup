package config

// UploadConfig holds configuration about how the backup will be uploaded.
type UploadConfig struct {
	// Host of the object storage API.
	Host string `validate:"required"`

	// KeyID for object storage API.
	KeyID string `validate:"required"`

	// SecretAccessKey for object storage API.
	SecretAccessKey string `validate:"required"`
}
