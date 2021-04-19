package db

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

// MigratePictures converts all indiewebdata already present in the database into local byte arrays (strings).
// This makes it possible to self-host author pictures. Run only after Migrate() in migrate-db.go.
func MigratePictures(cnf *common.Config, repo *MentionRepoBunt) {
	for _, domain := range cnf.AllowedWebmentionSources {
		all := repo.GetAll(domain)
		log.Info().Str("domain", domain).Int("mentions", len(all.Data)).Msg("migrate pictures: processing")
		for _, mention := range all.Data {
			if mention.Author.Picture == "" {
				log.Warn().Str("url", mention.Url).Msg("Mention without author picture, skipping")
				continue
			}

			savePicture(mention, repo, cnf)
		}
	}
}

// ChangeBaseUrl changes all base urls of pictures in the database.
// e.g. "http://localhost:1337/" to "https://jam.brainbaking.com/"
func ChangeBaseUrl(old, new string) {
	cnf := common.Configure()
	repo := NewMentionRepo(cnf)

	for _, domain := range cnf.AllowedWebmentionSources {
		for _, mention := range repo.GetAll(domain).Data {
			if mention.Author.Picture == "" {
				log.Warn().Str("url", mention.Url).Msg("Mention without author picture, skipping")
				continue
			}
			mention.Author.Picture = strings.ReplaceAll(mention.Author.Picture, old, new)
			repo.Save(mention.AsMention(), mention)
		}
	}
}

func savePicture(indieweb *mf.IndiewebData, repo *MentionRepoBunt, cnf *common.Config) {
	restClient := &rest.HttpClient{}
	picUrl := indieweb.Author.Picture
	log.Info().Str("oldurl", picUrl).Msg("About to cache picture")
	_, picData, err := restClient.GetBody(picUrl)
	if err != nil {
		log.Warn().Err(err).Str("url", picUrl).Msg("Unable to download author picture. Ignoring.")
		return
	}
	srcDomain := rest.Domain(indieweb.Source)
	_, dberr := repo.SavePicture(picData, srcDomain)
	if dberr != nil {
		log.Warn().Err(err).Str("url", picUrl).Msg("Unable to save downloaded author picture. Ignoring.")
		return
	}

	indieweb.Author.Picture = fmt.Sprintf("/pictures/%s", srcDomain)
	_, serr := repo.Save(indieweb.AsMention(), indieweb)
	if serr != nil {
		log.Fatal().Err(serr).Msg("Unable to update wm?")
	}
	log.Info().Str("oldurl", picUrl).Str("newurl", indieweb.Author.Picture).Msg("Picture saved!")
}
