package logging

import (
	"fmt"
	"gobit-demo/internal/config"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(cfg *config.AppConfig) {
	level, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		panic(fmt.Errorf("invalid log level: %w", err))
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
