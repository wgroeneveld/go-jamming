package main

import (
	"brainbaking.com/go-jamming/db"
	"flag"
	"os"

	"brainbaking.com/go-jamming/app"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	verboseFlag := flag.Bool("verbose", false, "Verbose mode (pretty print log, debug level)")
	migrateFlag := flag.Bool("migrate", false, "Run migration scripts for the DB and exit.")
	flag.Parse()

	// logs by default to Stderr (/var/log/syslog). Rolling files possible via lumberjack.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *verboseFlag || *migrateFlag {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if *migrateFlag {
		migrate()
		os.Exit(0)
	}

	log.Debug().Msg("Let's a go!")
	app.Start()
}

func migrate() {
	log.Info().Msg("Starting db migration...")
	db.Migrate()
	log.Info().Msg("Migration ended, exiting.")
}
