package backup

import (
	"io"
)

// Backuper performs the action of backing up a file.
type Backuper interface {
	// Backup files to w.
	Backup(w io.Writer) error
}
