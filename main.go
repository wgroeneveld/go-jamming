package main

import (
	"brainbaking.com/go-jamming/app/external"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
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
	blacklist := flag.String("blacklist", "", "Blacklist a domain name (also cleans spam from DB)")
	importFile := flag.String("import", "", "Import mentions from an external source (i.e. webmention.io)")
	flag.Parse()
	blacklisting := len(*blacklist) > 1
	importing := len(*importFile) > 1

	// logs by default to Stderr (/var/log/syslog). Rolling files possible via lumberjack.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *verboseFlag || *migrateFlag || blacklisting || importing {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if *migrateFlag {
		migrate()
		os.Exit(0)
	}

	if blacklisting {
		blacklistDomain(*blacklist)
		os.Exit(0)
	}

	if importing {
		importWebmentionFile(*importFile)
		os.Exit(0)
	}

	log.Debug().Msg("Let's a go!")
	app.Start()
}

func importWebmentionFile(file string) {
	log.Info().Str("file", file).Msg("Starting import...")

	config := common.Configure()
	bootstrapper := external.ImportBootstrapper{
		RestClient: &rest.HttpClient{},
		Conf:       config,
		Repo:       db.NewMentionRepo(config),
	}
	bootstrapper.Import(file)
}

func blacklistDomain(domain string) {
	log.Info().Str("domain", domain).Msg("Blacklisting...")
	config := common.Configure()
	config.AddToBlacklist(domain)

	repo := db.NewMentionRepo(config)
	for _, domain := range config.AllowedWebmentionSources {
		repo.CleanupSpam(domain, config.Blacklist)
	}

	log.Info().Msg("Blacklist done, exiting.")
}

func migrate() {
	log.Info().Msg("Starting db migration...")
	db.Migrate()
	log.Info().Msg("Migration ended, exiting.")
}
