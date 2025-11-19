package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/git"
)

var (
	gitClient *git.Git
)

var rootCmd = &cobra.Command{
	Use:   "periscope",
	Short: "A CLI tool for recursively viewing file contents in a directory",
	Long:  "Periscope recursively scans a directory or Git repo and prints file contents.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cfg, err := config.Load()
	if err != nil {
		return
	}

	gitClient = git.New(&cfg.Git)
}
