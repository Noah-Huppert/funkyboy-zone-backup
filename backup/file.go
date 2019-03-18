package backup

import (
	"archive/tar"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Noah-Huppert/mountain-backup/config"

	"github.com/Noah-Huppert/golog"
	"github.com/deckarep/golang-set"
)

// FilesBackuper backs up files.
type FilesBackuper struct {
	// Cfg configures which files to backup.
	Cfg config.FileConfig
}

// globSet expands any shell globs in set items and returns the resulting set.
func globSet(in mapset.Set) (mapset.Set, error) {
	out := mapset.NewSet()

	inIt := in.Iterator()

	for iUntyped := range inIt.C {
		i := iUntyped.(string)
		expanded, err := filepath.Glob(i)
		if err != nil {
			return nil, fmt.Errorf("error expanding shell globs in item \"%s\": %s", i, err.Error())
		}

		for _, o := range expanded {
			out.Add(o)
		}
	}

	return out, nil
}

// absSet resolves paths in set items to be absolute and returns the resulting set
func absSet(in mapset.Set) (mapset.Set, error) {
	out := mapset.NewSet()

	inIt := in.Iterator()

	for iUntyped := range inIt.C {
		i := iUntyped.(string)

		o, err := filepath.Abs(i)
		if err != nil {
			return nil, fmt.Errorf("error resolving absolute path in item \"%s\": %s", i, err.Error())
		}

		out.Add(o)
	}

	return out, nil
}

// allFiles walks a directory and returns an array of all the files in the directory.
func allFiles(walkPaths mapset.Set) (mapset.Set, error) {
	files := mapset.NewSet()

	walkPathsIt := walkPaths.Iterator()

	for walkPathUntyped := range walkPathsIt.C {
		walkPath := walkPathUntyped.(string)

		err := filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			files.Add(path)
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("error walking \"%s\": %s", walkPath, err.Error())
		}
	}

	return files, nil
}

// Backup configured files.
func (b FilesBackuper) Backup(logger golog.Logger, w *tar.Writer) error {
	// {{{1 Convert config to sets
	// {{{2 Cfg.Files
	filesSet := mapset.NewSet()
	for _, f := range b.Cfg.Files {
		filesSet.Add(f)
	}

	// {{{2 Cfg.Exclude
	excludeSet := mapset.NewSet()
	for _, e := range b.Cfg.Exclude {
		excludeSet.Add(e)
	}

	// {{{1 Resolve paths in config to absolute paths
	// {{{2 Cfg.Files
	absFiles, err := absSet(filesSet)
	if err != nil {
		return fmt.Errorf("error resolving absolute paths in Files configuration field: %s", err.Error())
	}

	// {{{2 Cfg.Exclude
	absExclude, err := absSet(excludeSet)
	if err != nil {
		return fmt.Errorf("error resolving absolute paths in Exclude configuration field: %s", err.Error())
	}

	// {{{1 Expand any shell globs in Cfg
	// {{{2 Cfg.Files
	globedFiles, err := globSet(absFiles)
	if err != nil {
		return fmt.Errorf("error expanding shell globs in Files configuration field: %s", err.Error())
	}

	// {{{2 Cfg.Exclude
	globedExclude, err := globSet(absExclude)
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
	backupFiles := walkedFiles.Difference(walkedExclude)

	// {{{1 Write files
	backupFilesIt := backupFiles.Iterator()

	for fileUntyped := range backupFilesIt.C {
		file := fileUntyped.(string)

		// {{{2 Write tar header
		// {{{3 Get file info
		fileInfo, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("error stat-ing \"%s\": %s", file, err.Error())
		}

		// {{{3 Write header
		err = w.WriteHeader(&tar.Header{
			Name: file,
			Mode: int64(fileInfo.Mode().Perm()),
			Size: fileInfo.Size(),
		})

		// {{{2 Write file body
		// {{{3 Open file
		fileReader, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("error opening \"%s\" for reading: %s", file, err.Error())
		}

		// {{{3 Read file body
		body, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return fmt.Errorf("error reading \"%s\" file contents: %s", file, err.Error())
		}

		// {{{3 Write to tar
		if _, err = w.Write(body); err != nil {
			return fmt.Errorf("error writing \"%s\" to tar file: %s", file, err.Error())
		}

		logger.Info(file)
	}

	return nil
}
