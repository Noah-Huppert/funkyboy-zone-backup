package config

// FileConfig holds configuration about standard files to backup
type FileConfig struct {
	// Files to backup. Can include shell globs.
	Files []string `validate:"required"`

	// Exclude holds names of files / directories to exclude from backup.
	// Can include shell globs.
	Exclude []string
}
