package config

// UploadConfig holds configuration about how the backup will be uploaded.
type UploadConfig struct {
	// Endpoint of the object storage API.
	Endpoint string `validate:"required"`

	// KeyID for object storage API.
	KeyID string `validate:"required"`

	// SecretAccessKey for object storage API.
	SecretAccessKey string `validate:"required"`

	// Bucket is the name of the bucket of which the backup will be uploaded.
	Bucket string `validate:"required"`

	// Format of backup file name without extension. Strftime symbols may be used.
	Format string `default:"backup-%Y-%m-%d-%H:%M:%S"`
}
