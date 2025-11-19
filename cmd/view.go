package cmd

import (
	"github.com/spf13/cobra"
)

var (
	copyToClipboard bool
	stripComments   bool

	ignorePath    []string
	ignoreContent []string
)

var viewCmd = &cobra.Command{
	Use:   "view [directory]",
	Short: "Recursively print contents of all files in a directory",
	Long:  "The view command recursively scans the given directory and prints the contents of every file it finds. If no directory is provided, the current working directory is used.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viewDir := "."
		if len(args) > 0 {
			viewDir = args[0]
		}

		opts := ViewOptions{
			Dir:             viewDir,
			CopyToClipboard: copyToClipboard,
			StripComments:   stripComments,
			IgnorePath:      ignorePath,
			IgnoreContent:   ignoreContent,
		}

		return RunView(opts)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.Flags().BoolVarP(
		&copyToClipboard,
		"copy",
		"c",
		false,
		"copy output to clipboard instead of printing",
	)

	viewCmd.Flags().BoolVarP(
		&stripComments,
		"strip-comments",
		"z",
		false,
		"strip comment lines (#, //, --) from file contents",
	)

	viewCmd.Flags().StringSliceVarP(
		&ignorePath,
		"ignore-path",
		"i",
		nil,
		"ignore files and directories whose path or name matches this glob-like pattern (can be repeated)",
	)

	viewCmd.Flags().StringSliceVarP(
		&ignoreContent,
		"ignore-content",
		"I",
		nil,
		"ignore files that contain a line matching this glob-like pattern (can be repeated)",
	)
}
