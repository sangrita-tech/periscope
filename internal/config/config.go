package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Git    Git    `yaml:"git"`
	Logger Logger `yaml:"logger"`
}

type Git struct {
	CacheRoot  string        `yaml:"cacheRoot"`
	MaxRepoAge time.Duration `yaml:"maxRepoAge"`
}

type Logger struct {
	Level      string            `yaml:"level"`
	DevMode    bool              `yaml:"devMode"`
	BaseFields map[string]string `yaml:"baseFields"`
}

func ReadConfigFile(path string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return &cfg, nil
}
