package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	contentbuilder "github.com/sangrita-tech/periscope/internal/content_builder"
	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/preprocess"
	"github.com/sangrita-tech/periscope/internal/scanner"
	"github.com/spf13/cobra"
)

var viewDirCmd = &cobra.Command{
	Use:   "dir [path]",
	Short: "Recursively print contents of all files in a local directory",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) > 0 {
			root = args[0]
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

		_, err := fmt.Fprint(os.Stdout, result)
		return err
	},
}

func init() {
	viewCmd.AddCommand(viewDirCmd)
}
