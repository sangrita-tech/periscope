package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Git Git `yaml:"git"`
}

type Git struct {
	CacheRoot  string        `yaml:"cache_root" env-default:"periscope"`
	MaxRepoAge time.Duration `yaml:"max_repo_age"  env-default:"336h"`
}

func Load() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
