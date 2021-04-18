package main

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

func mai() {
	cnf := common.Configure()
	os.Remove(cnf.Connection)
	repo := db.NewMentionRepo(cnf)

	log.Info().Str("dbconfig", cnf.Connection).Msg("Starting migration...")
	for _, domain := range cnf.AllowedWebmentionSources {
		fmt.Printf("Processing domain %s...\n", domain)
		entries, err := os.ReadDir(fmt.Sprintf("%s/%s", cnf.DataPath, domain))
		if err != nil {
			log.Fatal().Err(err).Msg("Error while reading import path")
		}

		for _, file := range entries {
			filename := fmt.Sprintf("%s/%s/%s", cnf.DataPath, domain, file.Name())
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal().Str("file", filename).Err(err).Msg("Error while reading file")
			}

			var indiewebData mf.IndiewebData
			json.Unmarshal(data, &indiewebData)
			mention := mf.Mention{
				Source: indiewebData.Source,
				Target: indiewebData.Target,
			}

			log.Info().Stringer("wm", mention).Str("file", filename).Msg("Re-saving entry")
			repo.Save(mention, &indiewebData)
		}
	}

	log.Info().Str("dbconfig", cnf.Connection).Msg("Checking for since files...")
	for _, domain := range cnf.AllowedWebmentionSources {
		since, err := ioutil.ReadFile(fmt.Sprintf("%s/%s-since.txt", cnf.DataPath, domain))
		if err != nil {
			log.Warn().Str("domain", domain).Msg("No since found, skipping")
			continue
		}

		log.Info().Str("domain", domain).Str("since", string(since)).Msg("Saving since")
		repo.UpdateSince(domain, common.IsoToTime(string(since)))
	}

	log.Info().Str("dbconfig", cnf.Connection).Msg("Done! Check db")
}
