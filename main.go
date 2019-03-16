package main

import (
	"github.com/Noah-Huppert/funkyboy-zone-backup/config"

	"github.com/Noah-Huppert/goconf"
	"github.com/Noah-Huppert/golog"
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

	logger.Debugf("%#v", cfg)
}
