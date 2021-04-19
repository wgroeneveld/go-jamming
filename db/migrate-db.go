package db

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

const (
	dataPath = "data" // decoupled from config, change if needed
)

// MigrateDataFiles migrates from data/[domain]/md5hash.json files to the new key/value db.
// This is only needed if you've run go-jamming before the db migration.
func MigrateDataFiles(cnf *common.Config, repo *MentionRepoBunt) {
	for _, domain := range cnf.AllowedWebmentionSources {
		log.Info().Str("domain", domain).Msg("MigrateDataFiles: processing")
		entries, err := os.ReadDir(fmt.Sprintf("%s/%s", dataPath, domain))
		if err != nil {
			log.Warn().Err(err).Msg("Error while reading import path - migration could be already done...")
			continue
		}

		for _, file := range entries {
			filename := fmt.Sprintf("%s/%s/%s", dataPath, domain, file.Name())
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal().Str("file", filename).Err(err).Msg("Error while reading file")
			}

			var indiewebData mf.IndiewebData
			json.Unmarshal(data, &indiewebData)
			mention := indiewebData.AsMention()

			log.Info().Stringer("wm", mention).Str("file", filename).Msg("Re-saving entry")
			repo.Save(mention, &indiewebData)
		}
	}

	log.Info().Str("dbconfig", cnf.ConString).Msg("Checking for since files...")
	for _, domain := range cnf.AllowedWebmentionSources {
		since, err := ioutil.ReadFile(fmt.Sprintf("%s/%s-since.txt", dataPath, domain))
		if err != nil {
			log.Warn().Str("domain", domain).Msg("No since found, skipping")
			continue
		}

		log.Info().Str("domain", domain).Str("since", string(since)).Msg("Saving since")
		repo.UpdateSince(domain, common.IsoToTime(string(since)))
	}
}
