package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/output"
	"github.com/sangrita-tech/periscope/internal/scanner"
	"github.com/sangrita-tech/periscope/internal/transformer"
	"github.com/sangrita-tech/periscope/internal/transformer/transformers"

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

		if err := scanner.WalkProcessedFiles(
			root,
			m,
			pipeline,
			agg.HandleFile,
		); err != nil {
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
	viewCmd.AddCommand(viewDirCmd)
}
