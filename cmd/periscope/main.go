package main

import (
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/logger"
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

	log.Info().Msg("Periscope stopped")
}
