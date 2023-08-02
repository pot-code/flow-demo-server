package logging

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(level string) {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(fmt.Errorf("invalid log level: %w", err))
	}
	zerolog.SetGlobalLevel(l)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
