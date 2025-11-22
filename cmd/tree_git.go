package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/git"
	"github.com/sangrita-tech/periscope/internal/scanner"
	"github.com/spf13/cobra"
)

var treeGitCmd = &cobra.Command{
	Use:   "git [repo]",
	Short: "Clone repo and print its tree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		g := git.New(&cfg.Git)
		root, err := g.Fetch(repo, gitBranch)
		if err != nil {
			return fmt.Errorf("failed to clone repo: %w", err)
		}

		pathM := buildPathMatcher()

		absRoot, err := filepath.Abs(root)
		if err != nil {
			absRoot = root
		}

		fmt.Println(filepath.Base(absRoot))

		s := scanner.New(absRoot, pathM, makeTreeHandlers())
		return s.Walk()
	},
}

func init() {
	treeCmd.AddCommand(treeGitCmd)

	treeGitCmd.Flags().StringVarP(
		&gitBranch,
		"branch",
		"b",
		"",
		"Git branch to use (default: repository default branch)",
	)
}
