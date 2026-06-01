package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/logger"
	"github.com/sangrita-tech/periscope/internal/output"
	"github.com/sangrita-tech/periscope/internal/render"
	"github.com/sangrita-tech/periscope/internal/source"
	"github.com/sangrita-tech/periscope/internal/walker"
)

var Version = "dev"

func main() {
	configPath := flag.String("config", "", "YML config")
	flag.Parse()

	if *configPath == "" {
		log.Fatal().Msg("Config path is required")
	}

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatal().Err(err).Str("config_path", *configPath).Msg("Failed to read config")
	}

	l := logger.New(&cfg.Logger)
	log.Logger = l

	log.Info().Msg("Starting periscope")

	s, err := source.ResolveSource(".")

	walker := walker.New()
	entries, err := walker.Walk(s)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to walk directory")
	}

	r := render.NewContentRenderer()
	out := r.Render(s, entries)

	w := output.NewTerminalWriter(os.Stdout)
	output.NewClipboardWriter()
	err = w.Write(out)

	log.Info().Int("num_entries", len(entries)).Msg("Finished walking directory")

	log.Info().Msg("Periscope stopped")
}
