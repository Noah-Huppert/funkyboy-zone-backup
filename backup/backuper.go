package backup

import (
	"archive/tar"
)

// Backuper performs the action of backing up a file.
type Backuper interface {
	// Backup files to w.
	Backup(w *tar.Writer) error
}
