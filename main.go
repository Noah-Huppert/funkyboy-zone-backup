package main

import (
	"archive/tar"
	"fmt"
	"os"
	"time"

	"github.com/Noah-Huppert/mountain-backup/backup"
	"github.com/Noah-Huppert/mountain-backup/config"

	"github.com/Noah-Huppert/goconf"
	"github.com/Noah-Huppert/golog"
	"github.com/jehiah/go-strftime"
)

func main() {
	// {{{1 Setup log
	logger := golog.NewStdLogger("backup")

	// {{{1 Load configuration
	cfgLoader := goconf.NewDefaultLoader()

	cfgLoader.AddConfigPath("./*.toml")
	cfgLoader.AddConfigPath("/etc/funkyboy-zone-backup/*.toml")

	cfg := config.Config{}
	if err := cfgLoader.Load(&cfg); err != nil {
		logger.Fatalf("error loading configuration: %s", err.Error())
	}

	// {{{1 Tar file
	// {{{2 Open tar file
	fName := strftime.Format(cfg.Upload.Format, time.Now())
	tarFPath := fmt.Sprintf("/var/tmp/%s.tar", fName)

	tarF, err := os.Create(tarFPath)
	defer func() {
		if err = tarF.Close(); err != nil {
			logger.Fatalf("error closing tar file \"%s\": %s",
				tarFPath, err.Error())
		}
	}()
	defer func() {
		if err = os.Remove(tarFPath); err != nil {
			logger.Fatalf("error removing tar file \"%s\": %s",
				tarFPath, err.Error())
		}
	}()

	if err != nil {
		logger.Fatalf("error creating tar file \"%s\": %s", tarFPath,
			err.Error())
	}

	// {{{2 Create tar writer
	tarW := tar.NewWriter(tarF)

	defer func() {
		if err := tarW.Close(); err != nil {
			logger.Fatalf("error closing tar writer \"%s\": %s", tarFPath, err.Error())
		}
	}()

	// {{{1 Perform straight forward backup of files
	for key, c := range cfg.Files {
		backuperLogger := logger.GetChild(fmt.Sprintf("File.%s", key))

		backuperLogger.Infof("backing up File.%s", key)

		b := backup.FilesBackuper{
			Cfg: c,
		}

		if err = b.Backup(backuperLogger, tarW); err != nil {
			logger.Fatalf("error running file backup for \"%s\": %s", key, err.Error())
		}
	}

	logger.Infof("backup: %s", tarFPath)
}
