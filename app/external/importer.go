package external

import (
	"brainbaking.com/go-jamming/app/mf"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
)

type Importer interface {
	TryImport(data []byte) ([]*mf.IndiewebData, error)
}

func Import(file string) {
	log.Info().Str("file", file).Msg("Starting import...")

	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Err(err).Msg("Unable to read file")
		return
	}

	importers := []Importer{
		&WebmentionIOImporter{},
	}

	var convertedData []*mf.IndiewebData

	for _, i := range importers {
		convertedData, err = i.TryImport(bytes)
		if err != nil {
			log.Warn().Str("importType", reflect.TypeOf(i).String()).Msg("Importer failed: ")
		} else {
			break
		}
	}

	if convertedData == nil {
		log.Fatal().Msg("No suitable importer found for data, aborting import!")
		return
	}

	// TODO store author pictures locally (and mutate wm for local URL)
	// TODO strip content + trim?
	// TODO save converted data in db
	// TODO whitelist domains?
}
