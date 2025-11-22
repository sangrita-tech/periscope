package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/sangrita-tech/periscope/internal/scanner"
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
	treeCmd.AddCommand(treeDirCmd)
}
