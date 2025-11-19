package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
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

		info, err := os.Stat(viewDir)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("directory %q does not exist", viewDir)
			}
			if os.IsPermission(err) {
				return fmt.Errorf("permission denied for directory %q", viewDir)
			}
			return fmt.Errorf("failed to access directory %q: %v", viewDir, err)
		}

		if !info.IsDir() {
			return fmt.Errorf("%q is not a directory", viewDir)
		}

		var out strings.Builder

		write := func(s string) {
			if copyToClipboard {
				out.WriteString(s)
			}
			if !copyToClipboard {
				fmt.Print(s)
			}
		}

		logErr := func(action, filePath string, err error) {
			abs, _ := filepath.Abs(filePath)
			msg := fmt.Sprintf("%s - %s (%v)\n", action, abs, err)
			if copyToClipboard {
				out.WriteString(msg)
			}
			fmt.Fprint(os.Stderr, msg)
		}

		// матчим и по полному пути, и по basename
		matchesAny := func(patterns []string, s string) bool {
			base := filepath.Base(s)

			for _, pat := range patterns {
				if pat == "" {
					continue
				}

				ok1, err1 := filepath.Match(pat, s)
				if err1 == nil && ok1 {
					return true
				}

				// отдельно проверяем basename, чтобы ".git" матчило директорию ".git"
				if base != s {
					ok2, err2 := filepath.Match(pat, base)
					if err2 == nil && ok2 {
						return true
					}
				}
			}
			return false
		}

		firstFile := true

		walkErr := filepath.WalkDir(viewDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				logErr("access error", path, err)
				return nil
			}

			// СНАЧАЛА проверяем игнор по пути (и для файлов, и для директорий)
			if matchesAny(ignorePath, path) {
				if d.IsDir() {
					// не заходить внутрь этой директории
					return fs.SkipDir
				}
				// просто пропускаем файл
				return nil
			}

			// дальше нас интересуют только файлы
			if d.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				logErr("open failed", path, err)
				return nil
			}
			defer f.Close()

			absPath, err := filepath.Abs(path)
			if err != nil {
				logErr("resolve path failed", path, err)
				return nil
			}

			scanner := bufio.NewScanner(f)

			var fileBuf strings.Builder
			skipFile := false

			for scanner.Scan() {
				line := scanner.Text()

				// если содержимое попадает под ignoreContent — выкидываем файл целиком
				if matchesAny(ignoreContent, line) {
					skipFile = true
					break
				}

				if stripComments {
					trimmed := strings.TrimSpace(line)
					if strings.HasPrefix(trimmed, "#") ||
						strings.HasPrefix(trimmed, "//") ||
						strings.HasPrefix(trimmed, "--") {
						continue
					}
				}

				fileBuf.WriteString(line + "\n")
			}

			if err := scanner.Err(); err != nil {
				logErr("read failed", path, err)
				return nil
			}

			if skipFile {
				return nil
			}

			if firstFile {
				write(fmt.Sprintf("[FILE] %s\n\n", absPath))
				firstFile = false
			} else {
				write(fmt.Sprintf("\n[FILE] %s\n\n", absPath))
			}

			write(fileBuf.String())

			return nil
		})

		if walkErr != nil {
			return walkErr
		}

		if copyToClipboard {
			if err := clipboard.WriteAll(out.String()); err != nil {
				return fmt.Errorf("failed to copy to clipboard: %v", err)
			}
		}

		return nil
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
