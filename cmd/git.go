package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/git"
	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/output"
	"github.com/sangrita-tech/periscope/internal/scanner"
	"github.com/sangrita-tech/periscope/internal/transformer"
	"github.com/sangrita-tech/periscope/internal/transformer/transformers"

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
		root, err := g.CloneRepo(repo, gitBranch)
		if err != nil {
			return fmt.Errorf("failed to clone repo: %w", err)
		}

		m := matcher.New(ignorePaths, ignoreContents)

		pipeline := transformer.New().
			Add(transformers.CollapseEmptyLines())

		if stripComments {
			pipeline.Add(transformers.StripComments())
		}

		if maskURL {
			pipeline.Add(transformers.MaskURL())
		}

		agg := output.NewAggregator(copyToClipboard, os.Stdout)

		err = scanner.WalkProcessedFiles(
			root,
			m,
			pipeline,
			agg.HandleFile,
		)
		if err != nil {
			return err
		}

		if copyToClipboard {
			if err := clipboard.WriteAll(agg.Result()); err != nil {
				return fmt.Errorf("failed to copy to clipboard: %w", err)
			}
		}

		return nil
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
