package backup

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/Noah-Huppert/mountain-backup/config"
)

// FileBackuper backs up files.
type FileBackuper struct {
	// Cfg configures which files to backup.
	Cfg config.FileConfig
}

// globArray expands any shell globs in array items and returns the resulting array.
func globArray(in []string) ([]string, error) {
	out := []string{}

	for _, i := range in {
		expanded, err := filepath.Glob(i)
		if err != nil {
			return nil, fmt.Errorf("error expanding shell globs in item \"%s\": %s", i, err.Error())
		}

		for _, o := range expanded {
			out = append(out, o)
		}
	}

	return out
}

// allFiles walks a directory and returns an array of all the files in the directory.
func allFiles(dir string) ([]string, error) {
	// TODO: Walk and return array
}

// Backup configured files.
func (b FileBackuper) Backup(w io.Writer) error {
	// {{{1 Expand any shell globs in Cfg
	// {{{2 Cfg.Files
	globedFiles, err := globArray(b.Cfg.Files)
	if err != nil {
		return fmt.Errorf("error expanding shell globs in Files configuration field: %s", err.Error())
	}

	// {{{2 Cfg.Exclude
	globedExclude, err := globArray(b.Cfg.Exclude)
	if err != nil {
		return fmt.Errorf("error expanding shell globs in Exclude configuration field: %s", err.Error())
	}

	// {{{1 Walk any directories to create a complete list of files
	// {{{2 Cfg.Files
	walkedFiles, err := allFiles(globedFiles)
	if err != nil {
		return fmt.Errorf("error creating list of all files in Files configuration field: %s", err.Error())
	}

	// {{{2 Cfg.Exclude
	walkedExclude, err := allFiles(globedExclude)
	if err != nil {
		return fmt.Errorf("error creating list of all files in Exclude configuration field: %s", err.Error())
	}

	// TODO: Remove excluded files from list
	// TODO: Write files to w
}
