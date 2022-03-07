package config

import (
	"strings"

	"github.com/kkyr/fig"
	"github.com/pkg/errors"
)

type Config struct {
	BindAddr string `fig:"bind_addr" default:""`
	Port     int    `fig:"port" default:"3690"`

	YoutubeDLCommand string `fig:"youtubedl_command" default:"python3 /usr/local/bin/youtube-dl"`

	DataDir string `fig:"data_dir" validate:"required"`

	PlexBaseURL string `fig:"plex_base_url"`
	PlexToken   string `fig:"plex_token"`

	LibraryDir   string `fig:"library_dir" validate:"required"`
	LibraryOwner string `fig:"library_owner"`
}

var defaultConfig = Config{
	Port: 3690,
}

func Read(cfgFile string) (cfg Config, err error) {
	err = fig.Load(&cfg, fig.UseEnv("BROADTAIL"), fig.File(cfgFile))
	if errors.Is(err, fig.ErrFileNotFound) {
		return defaultConfig, err
	}
	return cfg, nil
}

func (cfg Config) YoutubeDLCommandAsSlice() []string {
	return strings.Split(cfg.YoutubeDLCommand, " ")
}
