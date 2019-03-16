package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Noah-Huppert/funkyboy-zone-backup/config"

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

	// {{{1 Open tar file
	fName := fmt.Sprintf("backup-%s", strftime.Format("%Y-%m-%d-%H:%M:%S",
		time.Now()))
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
}
