package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
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

		result, err := runViewScan(root)
		if err != nil {
			return runtimeErr(cmd, err)
		}

		if copyToClipboard {
			if err := clipboard.WriteAll(result); err != nil {
				return runtimeErr(cmd, fmt.Errorf("failed to copy to clipboard: %w", err))
			}
			return nil
		}

		if _, err := fmt.Fprint(os.Stdout, result); err != nil {
			return runtimeErr(cmd, err)
		}
		return nil
	},
}

func init() {
	viewCmd.AddCommand(viewDirCmd)
}
