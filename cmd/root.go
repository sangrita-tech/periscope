package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "periscope",
	Short: "A CLI tool for recursively viewing file contents in a directory",
	Long:  "Periscope is a small cross-platform CLI utility that recursively scans a directory and prints the full contents of every file it finds. It is useful for quick inspection, debugging, and exploring unknown file structures.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
