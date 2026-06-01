package config

import (
	"fmt"
	"os"

	"github.com/sangrita-tech/periscope/internal/model"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Ignore  []string            `yaml:"ignore"`
	Replace []model.Replacement `yaml:"replace"`
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
