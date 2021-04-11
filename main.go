package main

import (
	"flag"
	"os"

	"brainbaking.com/go-jamming/app"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	verboseFlag := flag.Bool("verbose", false, "Verbose mode (pretty print log, debug level)")
	flag.Parse()

	// logs by default to Stderr (/var/log/syslog). Rolling files possible via lumberjack.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *verboseFlag == true {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Debug().Msg("Let's a go!")
	app.Start()
}
