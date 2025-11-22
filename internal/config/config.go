package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const envConfigPath = "PERISCOPE_CONFIG"

type Config struct {
	Git Git `yaml:"git"`
}

type Git struct {
	CacheRoot  string        `yaml:"cache_root" env-default:"periscope"`
	MaxRepoAge time.Duration `yaml:"max_repo_age"  env-default:"336h"`
}

func Load() (*Config, error) {
	var cfg Config

	cfgFile := os.Getenv(envConfigPath)
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			cfgFile = filepath.Join(home, ".periscope.yaml")
		}
	}

	if cfgFile == "" {
		_ = cleanenv.ReadEnv(&cfg)
		return &cfg, nil
	}

	if _, err := os.Stat(cfgFile); errors.Is(err, os.ErrNotExist) {
		_ = cleanenv.ReadEnv(&cfg)
		return &cfg, nil
	}

	if err := cleanenv.ReadConfig(cfgFile, &cfg); err != nil {
		return nil, err
	}

	_ = cleanenv.ReadEnv(&cfg)
	return &cfg, nil
}
