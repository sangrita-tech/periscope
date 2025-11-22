package cmd

import (
	"github.com/sangrita-tech/periscope/internal/ui"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"

	copyToClipboard bool
	stripComments   bool
	maskURL         bool
	ignorePaths     []string
	ignoreContents  []string

	log = ui.New()
)

var rootCmd = &cobra.Command{
	Use:   "periscope",
	Short: "A CLI tool for recursively viewing file contents in a directory",
	Long:  "Periscope recursively scans a directory or Git repo and prints file contents.",
}

func Execute() {
	_ = rootCmd.Execute()
}

func init() {
	rootCmd.Version = Version
}

func runtimeErr(cmd *cobra.Command, err error) error {
	if err == nil {
		return nil
	}
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	log.Error("%v", err)
	return err
}
