package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/git"
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

		result, err := runTreeScan(root)
		if err != nil {
			return err
		}

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
	treeCmd.AddCommand(treeGitCmd)

	treeGitCmd.Flags().StringVarP(
		&gitBranch,
		"branch",
		"b",
		"",
		"Git branch to use (default: repository default branch)",
	)
}
