package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/git"
	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/transformer"

	"github.com/spf13/cobra"
)

var (
	gitBranch string
)

var viewGitCmd = &cobra.Command{
	Use:   "git [repo]",
	Short: "Clone and view contents of a Git repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		g := git.New(&cfg.Git)

		root, err := g.CloneRepo(repo, gitBranch)
		if err != nil {
			return fmt.Errorf("failed to clone repo: %w", err)
		}

		m := matcher.New(ignorePaths, ignoreContents)

		pipeline := transformer.NewPipeline()
		if stripComments {
			pipeline.Use(transformer.StripComments())
		}
		pipeline.Use(transformer.CollapseEmptyLines())

		var (
			buf       bytes.Buffer
			firstFile = true
		)

		err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				if m.ShouldIgnorePath(path) {
					return filepath.SkipDir
				}
				return nil
			}

			if m.ShouldIgnorePath(path) {
				return nil
			}

			contentBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if bytes.IndexByte(contentBytes, 0x00) != -1 {
				return nil
			}

			content := string(contentBytes)

			if m.ShouldIgnoreContent(content) {
				return nil
			}

			content, err = pipeline.Process(path, content)
			if err != nil {
				return err
			}

			if strings.TrimSpace(content) == "" {
				return nil
			}

			absPath, err := filepath.Abs(path)
			if err != nil {
				absPath = path
			}

			header := fmt.Sprintf("[FILE] %s\n\n", absPath)

			if !copyToClipboard {
				if !firstFile {
					fmt.Println()
				}

				fmt.Print(header)
				fmt.Print(content)
				if !strings.HasSuffix(content, "\n") {
					fmt.Println()
				}
			}

			if !firstFile {
				buf.WriteString("\n")
			}
			buf.WriteString(header)
			buf.WriteString(content)
			if !strings.HasSuffix(content, "\n") {
				buf.WriteString("\n")
			}

			firstFile = false
			return nil
		})

		if err != nil {
			return err
		}

		if copyToClipboard {
			if err := clipboard.WriteAll(buf.String()); err != nil {
				return fmt.Errorf("failed to copy to clipboard: %w", err)
			}
		}

		return nil
	},
}

func init() {
	viewCmd.AddCommand(viewGitCmd)

	viewGitCmd.Flags().StringVarP(
		&gitBranch,
		"branch",
		"b",
		"",
		"Git branch to view (default: repository default branch)",
	)
}
