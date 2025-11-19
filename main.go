package main

import (
	"log"

	"github.com/sangrita-tech/periscope/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
