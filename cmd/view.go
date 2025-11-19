package cmd

import (
	"github.com/spf13/cobra"
)

var (
	copyToClipboard bool
	stripComments   bool

	ignorePaths    []string
	ignoreContents []string
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Recursively print contents of files (local dir or git repo)",
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.PersistentFlags().BoolVarP(&copyToClipboard, "copy", "c", false, "copy result to clipboard")
	viewCmd.PersistentFlags().BoolVarP(&stripComments, "strip-comments", "z", false, "remove comment lines")

	viewCmd.PersistentFlags().StringSliceVarP(&ignorePaths, "ignore-path", "i", nil, "ignore files/dirs matching pattern")
	viewCmd.PersistentFlags().StringSliceVarP(&ignoreContents, "ignore-content", "I", nil, "ignore files whose content matches pattern")
}
