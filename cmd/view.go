package cmd

import (
	"bufio"
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	contentbuilder "github.com/sangrita-tech/periscope/internal/content_builder"
	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/preprocess"
	"github.com/sangrita-tech/periscope/internal/scanner"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Recursively print contents of files (local dir or git repo)",
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.PersistentFlags().BoolVarP(&maskURL, "mask-url", "m", false, "mask urls in files")
	viewCmd.PersistentFlags().BoolVarP(&copyToClipboard, "copy", "c", false, "copy result to clipboard")
	viewCmd.PersistentFlags().BoolVarP(&stripComments, "strip-comments", "z", false, "remove comment lines")

	viewCmd.PersistentFlags().StringSliceVarP(&ignorePaths, "ignore-path", "i", nil, "ignore files/dirs matching pattern")
	viewCmd.PersistentFlags().StringSliceVarP(&ignoreContents, "ignore-content", "I", nil, "ignore files whose content matches pattern")
}

func buildPathMatcher() *matcher.Matcher {
	if len(ignorePaths) == 0 {
		return nil
	}
	return matcher.New(ignorePaths)
}

func buildContentMatcher() *matcher.Matcher {
	if len(ignoreContents) == 0 {
		return nil
	}
	return matcher.New(ignoreContents)
}

func buildChain() *preprocess.Chain {
	ch := preprocess.New().AddCollapseEmptyLines()
	if stripComments {
		ch.AddStripComments()
	}
	if maskURL {
		ch.AddMaskURL()
	}
	return ch
}

func runViewScan(root string) (string, error) {
	pathM := buildPathMatcher()
	contentM := buildContentMatcher()
	chain := buildChain()
	builder := contentbuilder.New()

	handlers := scanner.Handlers{
		OnFile: func(path string, d fs.DirEntry, depth int, isLast bool) error {
			contentBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if bytes.IndexByte(contentBytes, 0x00) != -1 {
				return nil
			}

			content := string(contentBytes)

			if contentM != nil {
				sc := bufio.NewScanner(strings.NewReader(content))
				for sc.Scan() {
					if contentM.Match(sc.Text()) {
						return nil
					}
				}
				if err := sc.Err(); err != nil {
					return err
				}
			}

			processed, _, err := chain.Process(path, content)
			if err != nil {
				return err
			}
			content = processed

			if strings.TrimSpace(content) == "" {
				return nil
			}

			absPath, err := filepath.Abs(path)
			if err != nil {
				absPath = path
			}

			return builder.AddBlock("[FILE] "+absPath, content)
		},
	}

	s := scanner.New(root, pathM, handlers)
	if err := s.Walk(); err != nil {
		return "", err
	}

	return builder.Result(), nil
}
