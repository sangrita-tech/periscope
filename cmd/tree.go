package cmd

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/sangrita-tech/periscope/internal/scanner"
	"github.com/spf13/cobra"
)

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Print directory structure as a tree (local dir or git repo)",
}

func init() {
	rootCmd.AddCommand(treeCmd)

	treeCmd.PersistentFlags().BoolVarP(&copyToClipboard, "copy", "c", false, "copy result to clipboard")

	treeCmd.PersistentFlags().StringSliceVarP(&ignorePaths, "ignore-path", "i", nil, "ignore files/dirs matching pattern")
	treeCmd.PersistentFlags().StringSliceVarP(&ignoreContents, "ignore-content", "I", nil, "ignore files whose content matches pattern")
}

func makeTreeHandlers() scanner.Handlers {
	prefixes := []string{}

	onNode := func(d fs.DirEntry, depth int, isLast bool) {
		for len(prefixes) < depth {
			prefixes = append(prefixes, "")
		}

		prefix := strings.Join(prefixes[:depth-1], "")

		branch := "├─ "
		if isLast {
			branch = "└─ "
		}

		fmt.Println(prefix + branch + d.Name())

		if isLast {
			prefixes[depth-1] = "   "
		} else {
			prefixes[depth-1] = "│  "
		}
	}

	return scanner.Handlers{
		OnDir: func(path string, d fs.DirEntry, depth int, isLast bool) error {
			onNode(d, depth, isLast)
			return nil
		},
		OnFile: func(path string, d fs.DirEntry, depth int, isLast bool) error {
			onNode(d, depth, isLast)
			return nil
		},
	}
}
