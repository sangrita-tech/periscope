package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Ignore  []string      `yaml:"ignore"`
	Replace []Replacement `yaml:"replace"`
}

type Replacement struct {
	Pattern string `yaml:"pattern"`
	Value   string `yaml:"value"`
}

func (cfg *Config) ApplyReplacements(value string) string {
	if cfg == nil || len(cfg.Replace) == 0 {
		return value
	}

	for _, replacement := range cfg.Replace {
		if replacement.Pattern == "" {
			continue
		}

		value = strings.ReplaceAll(value, replacement.Pattern, replacement.Value)
	}

	return value
}

func ReadConfig(path string) (*Config, error) {
	path, err := expandPath(path)
	if err != nil {
		return nil, err
	}

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

func expandPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("empty config path")
	}

	if path == "~" || strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}

		if path == "~" {
			return homeDir, nil
		}

		path = filepath.Join(homeDir, path[2:])
	}

	return path, nil
}
