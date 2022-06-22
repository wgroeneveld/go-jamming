package external

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/notifier"
	"brainbaking.com/go-jamming/app/webmention/recv"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
)

type Importer interface {
	TryImport(data []byte) ([]*mf.IndiewebData, error)
}

type ImportBootstrapper struct {
	RestClient rest.Client
	Conf       *common.Config
	Repo       db.MentionRepo
}

func (ib *ImportBootstrapper) Import(file string) {
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
		if err != nil || len(convertedData) == 0 {
			log.Warn().Str("importType", reflect.TypeOf(i).String()).Msg("Importer failed or returned zero entries")
		} else {
			log.Info().Str("importType", reflect.TypeOf(i).String()).Msg("Suitable converter found!")
			break
		}
	}

	if convertedData == nil {
		log.Fatal().Msg("No suitable importer found for data, aborting import!")
		return
	}

	log.Info().Msg("Conversion succeeded, persisting to data storage...")
	recv := &recv.Receiver{
		RestClient: ib.RestClient,
		Conf:       ib.Conf,
		Repo:       ib.Repo,
		Notifier:   &notifier.StringNotifier{},
	}

	for _, wm := range convertedData {
		mention := mf.Mention{
			Source: wm.Source,
			Target: wm.Target,
		}
		ib.Conf.AddToWhitelist(mention.SourceDomain())

		recv.ProcessAuthorPicture(wm)
		recv.ProcessWhitelistedMention(mention, wm)
	}

	log.Info().Msg("All done, enjoy your go-jammed mentions!")
}
