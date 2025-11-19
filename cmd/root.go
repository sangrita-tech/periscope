package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "periscope",
	Short: "A CLI tool for recursively viewing file contents in a directory",
	Long:  "Periscope recursively scans a directory or Git repo and prints file contents.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Version = Version
}
