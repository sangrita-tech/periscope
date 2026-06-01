package main

import (
	"os"

	"github.com/sangrita-tech/periscope/internal/cli"
)

var Version = "dev"

func main() {
	if err := cli.Execute(Version); err != nil {
		os.Exit(1)
	}
}
