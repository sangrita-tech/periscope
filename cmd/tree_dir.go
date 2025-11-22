package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var treeDirCmd = &cobra.Command{
	Use:   "dir [path]",
	Short: "Print tree of a local directory",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) > 0 {
			root = args[0]
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
	treeCmd.AddCommand(treeDirCmd)
}
