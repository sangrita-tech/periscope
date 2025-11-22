package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/git"

	contentbuilder "github.com/sangrita-tech/periscope/internal/content_builder"
	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/preprocess"
	"github.com/sangrita-tech/periscope/internal/scanner"

	"github.com/spf13/cobra"
)

var (
	gitBranch string
)

var viewGitCmd = &cobra.Command{
	Use:   "git [repo]",
	Short: "Clone and view contents of a Git repository",
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

		var pathM *matcher.Matcher
		if len(ignorePaths) > 0 {
			pathM = matcher.New(ignorePaths)
		}

		var contentM *matcher.Matcher
		if len(ignoreContents) > 0 {
			contentM = matcher.New(ignoreContents)
		}

		chain := preprocess.New().
			AddCollapseEmptyLines()

		if stripComments {
			chain.AddStripComments()
		}

		if maskURL {
			chain.AddMaskURL()
		}

		builder := contentbuilder.New()

		s := scanner.New(root, pathM, contentM, chain, builder)
		if err := s.Walk(); err != nil {
			return err
		}

		result := builder.Result()

		if copyToClipboard {
			if err := clipboard.WriteAll(result); err != nil {
				return fmt.Errorf("failed to copy to clipboard: %w", err)
			}
			return nil
		}

		_, err = fmt.Fprint(os.Stdout, result)
		return err
	},
}

func init() {
	viewCmd.AddCommand(viewGitCmd)

	viewGitCmd.Flags().StringVarP(
		&gitBranch,
		"branch",
		"b",
		"",
		"Git branch to view (default: repository default branch)",
	)
}
