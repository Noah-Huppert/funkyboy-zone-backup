package backup

import (
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Noah-Huppert/mountain-backup/config"
)

// FilesBackuper backs up files.
type FilesBackuper struct {
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

	return out, nil
}

// allFiles walks a directory and returns an array of all the files in the directory.
func allFiles(walkPaths []string) ([]string, error) {
	files := []string{}

	for _, walkPath := range walkPaths {
		err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			files = append(files, path)

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("error walking \"%s\": %s", walkPath, err.Error())
		}
	}

	return files, nil
}

// Backup configured files.
func (b FilesBackuper) Backup(w *tar.Writer) error {
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

	// {{{1 Remove excluded files
	backupFiles := []string{}

	for _, f := range walkedFiles {
		excluded := false
		for _, e := range walkedExclude {
			if f == e {
				excluded = true
				break
			}
		}

		if excluded {
			break
		}

		backupFiles = append(backupFiles, f)
	}

	// {{{1 Write files
	fmt.Printf("%#v\n", backupFiles)
	return nil
}
