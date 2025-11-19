package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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

		return filepath.WalkDir(viewDir, func(path string, d os.DirEntry, err error) error {
			logErr := func(action, filePath string, err error) {
				abs, _ := filepath.Abs(filePath)
				fmt.Fprintf(os.Stderr, "%s - %s (%v)\n", action, abs, err)
			}

			if err != nil {
				logErr("access error", path, err)
				return nil
			}

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

			fmt.Printf("\n[FILE] %s\n\n", absPath)

			if _, err := io.Copy(os.Stdout, f); err != nil {
				logErr("read failed", path, err)
			}

			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
