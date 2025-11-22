package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sangrita-tech/periscope/internal/scanner"
	"github.com/sangrita-tech/periscope/internal/ui"
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

func makeTreeHandlers(buf *bytes.Buffer) scanner.Handlers {
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

		name := d.Name()
		if d.IsDir() {
			name = ui.DirStyle.Sprintf("📁 %s", name)
		}

		line := ui.Branch.Sprintf(prefix+branch) + name
		fmt.Fprintln(buf, line)

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

func runTreeScan(root string) (string, error) {
	pathM := buildPathMatcher()

	absRoot, err := filepath.Abs(root)
	if err != nil {
		absRoot = root
	}

	var buf bytes.Buffer

	fmt.Fprintln(&buf, ui.Title.Sprintf("📦 %s", filepath.Base(absRoot)))

	s := scanner.New(absRoot, pathM, makeTreeHandlers(&buf))
	if err := s.Walk(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
