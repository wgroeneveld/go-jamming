
package main

import (
	"os"

	"brainbaking.com/go-jamming/app"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    // TODO this should only be enabled in local mode. Fix with config?
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

    log.Debug().Msg("Let's a go!")
	app.Start()
}
