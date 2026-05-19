package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger Logger `yaml:"logger"`
}

type Logger struct {
	Level      string            `yaml:"level"`
	DevMode    bool              `yaml:"devMode"`
	BaseFields map[string]string `yaml:"baseFields"`
}

func (cfg *Config) applyDefaults() {
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}
	if cfg.Logger.BaseFields == nil {
		cfg.Logger.BaseFields = make(map[string]string)
	}
}

func ReadConfigFile(path string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg.applyDefaults()

	return &cfg, nil
}
