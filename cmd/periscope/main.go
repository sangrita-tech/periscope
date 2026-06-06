package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/ignore"
	"github.com/sangrita-tech/periscope/internal/output"
	"github.com/sangrita-tech/periscope/internal/render"
	"github.com/sangrita-tech/periscope/internal/replacement"
	"github.com/sangrita-tech/periscope/internal/source"
	"github.com/sangrita-tech/periscope/internal/walker"
	"github.com/spf13/cobra"
)

const configFileName = ".periscope.yml"

var Version = "dev"

func main() {
	var configPath string
	var ignorePatterns []string
	var treeOnly bool
	var copyToClipboard bool

	cmd := &cobra.Command{
		Use:           "periscope [path]",
		Short:         "Print a text snapshot of a local path or Git repository",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			target := "."
			if len(args) > 0 {
				target = args[0]
			}

			return run(
				target,
				configPath,
				ignorePatterns,
				treeOnly,
				copyToClipboard,
			)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "YAML config path")
	cmd.Flags().StringArrayVarP(&ignorePatterns, "ignore", "i", nil, "ignore file or directory pattern (repeatable)")
	cmd.Flags().BoolVarP(&treeOnly, "tree", "t", false, "print only the file tree")
	cmd.Flags().BoolVarP(&copyToClipboard, "copy", "c", false, "copy the result to the clipboard")

	cmd.Version = Version

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(
	target string,
	configPath string,
	ignorePatterns []string,
	treeOnly bool,
	copyToClipboard bool,
) error {

	cfg, err := initConfig(configPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	ignorePatterns = splitIgnorePatterns(ignorePatterns)

	allIgnorePatterns := append([]string{}, cfg.Ignore...)
	allIgnorePatterns = append(allIgnorePatterns, ignorePatterns...)

	matcher, err := ignore.NewMatcher(ignorePatterns)
	if err != nil {
		return fmt.Errorf("build ignore matcher: %w", err)
	}

	source, err := source.ResolveSource(target)
	if err != nil {
		return fmt.Errorf("resolve source: %w", err)
	}

	entries, err := walker.New(matcher).Walk(source)
	if err != nil {
		return fmt.Errorf("walk source: %w", err)
	}

	currentTime := time.Now()
	renderedEntries := ""
	if treeOnly {
		renderedEntries = render.RenderTree(source, entries, currentTime)
	} else {
		renderedEntries = render.RenderContent(source, entries, currentTime)
	}

	result := replacement.Apply(renderedEntries, cfg.Replace)

	if err := output.NewTerminalWriter(os.Stdout).Write(result); err != nil {
		return fmt.Errorf("write terminal output: %w", err)
	}

	if copyToClipboard {
		if err := output.NewClipboardWriter().Write(result); err != nil {
			return fmt.Errorf("copy output to clipboard: %w", err)
		}
	}

	return nil
}

func initConfig(configPath string) (*config.Config, error) {
	if configPath != "" {
		return config.ReadConfig(configPath)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve home directory: %w", err)
	}

	defaultConfigPath := filepath.Join(homeDir, configFileName)
	cfg, err := config.ReadConfig(defaultConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &config.Config{}, nil
		}
		return nil, err
	}

	return cfg, nil
}

func splitIgnorePatterns(patterns []string) []string {
	var result []string

	for _, pattern := range patterns {
		for _, part := range strings.Split(pattern, "|") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			result = append(result, part)
		}
	}

	return result
}
