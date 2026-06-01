package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/ignore"
	"github.com/sangrita-tech/periscope/internal/output"
	"github.com/sangrita-tech/periscope/internal/render"
	"github.com/sangrita-tech/periscope/internal/source"
	"github.com/sangrita-tech/periscope/internal/walker"
	"github.com/spf13/cobra"
)

const baseConfigPath = "~/.periscope.yml"

type options struct {
	configPath string
	ignore     []string
	tree       bool
	copy       bool
	version    string
}

func Execute(version string) error {
	opts := options{
		version: version,
	}

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

			return run(opts, target)
		},
	}

	cmd.Flags().StringVar(&opts.configPath, "config", "", "YAML config path")
	cmd.Flags().StringArrayVarP(&opts.ignore, "ignore", "i", nil, "ignore file or directory pattern (repeatable)")
	cmd.Flags().BoolVarP(&opts.tree, "tree", "t", false, "print only the file tree")
	cmd.Flags().BoolVarP(&opts.copy, "copy", "c", false, "copy the result to the clipboard")

	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	return nil
}

func run(opts options, target string) error {
	cfg, err := initConfig(opts.configPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	ignorePatterns := append([]string{}, cfg.Ignore...)
	ignorePatterns = append(ignorePatterns, opts.ignore...)

	ignoreMatcher, err := ignore.NewMatcher(ignorePatterns)
	if err != nil {
		return fmt.Errorf("build ignore matcher: %w", err)
	}

	s, err := source.ResolveSource(target)
	if err != nil {
		return fmt.Errorf("resolve source: %w", err)
	}

	w := walker.New(ignoreMatcher)

	entries, err := w.Walk(s)
	if err != nil {
		return fmt.Errorf("walk source: %w", err)
	}

	var r render.Renderer
	if opts.tree {
		r = render.NewTreeRenderer()
	} else {
		r = render.NewContentRenderer()
	}

	result := r.Render(s, entries)

	if err := output.NewTerminalWriter(os.Stdout).Write(result); err != nil {
		return fmt.Errorf("write terminal output: %w", err)
	}

	if opts.copy {
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

	cfg, err := config.ReadConfig(baseConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &config.Config{}, nil
		}

		return nil, err
	}

	return cfg, nil
}
