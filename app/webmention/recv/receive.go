package recv

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/pictures"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"willnorris.com/go/microformats"
)

type Receiver struct {
	RestClient rest.Client
	Conf       *common.Config
	Repo       db.MentionRepo
}

func (recv *Receiver) Receive(wm mf.Mention) {
	log.Info().Stringer("wm", wm).Msg("OK: looks valid")
	_, body, geterr := recv.RestClient.GetBody(wm.Source)

	if geterr != nil {
		log.Warn().Err(geterr).Msg("  ABORT: invalid url")
		recv.Repo.Delete(wm)
		return
	}

	recv.processSourceBody(body, wm)
}

func (recv *Receiver) processSourceBody(body string, wm mf.Mention) {
	if !strings.Contains(body, wm.Target) {
		log.Warn().Str("target", wm.Target).Msg("ABORT: no mention of target found in html src of source!")
		return
	}

	data := microformats.Parse(strings.NewReader(body), wm.SourceUrl())
	indieweb := recv.convertBodyToIndiewebData(body, wm, mf.HEntry(data))
	if indieweb.Author.Picture != "" {
		recv.saveAuthorPictureLocally(indieweb)
	}

	key, err := recv.Repo.Save(wm, indieweb)
	if err != nil {
		log.Error().Err(err).Stringer("wm", wm).Msg("processSourceBody: failed to save json to db")
	}
	log.Info().Str("key", key).Msg("OK: Webmention processed.")
}

func (recv *Receiver) convertBodyToIndiewebData(body string, wm mf.Mention, hEntry *microformats.Microformat) *mf.IndiewebData {
	if hEntry == nil {
		return recv.parseBodyAsNonIndiewebSite(body, wm)
	}
	return recv.parseBodyAsIndiewebSite(hEntry, wm)
}

// see https://github.com/willnorris/microformats/blob/main/microformats.go
func (recv *Receiver) parseBodyAsIndiewebSite(hEntry *microformats.Microformat, wm mf.Mention) *mf.IndiewebData {
	return &mf.IndiewebData{
		Name: mf.Str(hEntry, "name"),
		Author: mf.IndiewebAuthor{
			Name:    mf.DetermineAuthorName(hEntry),
			Picture: mf.Str(mf.Prop(hEntry, "author"), "photo"),
		},
		Content:      mf.Content(hEntry),
		Url:          mf.Url(hEntry, wm.Source),
		Published:    mf.Published(hEntry, recv.Conf.UtcOffset),
		Source:       wm.Source,
		Target:       wm.Target,
		IndiewebType: mf.Type(hEntry),
	}
}

var (
	titleRegexp = regexp.MustCompile(`<title>(.*?)<\/title>`)
)

func (recv *Receiver) parseBodyAsNonIndiewebSite(body string, wm mf.Mention) *mf.IndiewebData {
	title := nonIndiewebTitle(body, wm)
	return &mf.IndiewebData{
		Author: mf.IndiewebAuthor{
			Name: wm.Source,
		},
		Name:         title,
		Content:      title,
		Published:    mf.PublishedNow(recv.Conf.UtcOffset),
		Url:          wm.Source,
		IndiewebType: mf.TypeMention,
		Source:       wm.Source,
		Target:       wm.Target,
	}
}

// saveAuthorPictureLocally tries to download the author picture.
// If it succeeds, it alters the picture path to a local /pictures/x one.
// If it fails, it falls back to a local default image.
// This *should* also validate image byte headers, like https://stackoverflow.com/questions/670546/determine-if-file-is-an-image
func (recv *Receiver) saveAuthorPictureLocally(indieweb *mf.IndiewebData) {
	_, picData, err := recv.RestClient.GetBody(indieweb.Author.Picture)
	if err != nil {
		log.Warn().Err(err).Str("url", indieweb.Author.Picture).Msg("Unable to download author picture. Reverting to anonymous.")
		indieweb.Author.Picture = fmt.Sprintf("/pictures/%s", pictures.Anonymous)
		return
	}
	srcDomain := rest.Domain(indieweb.Source)
	_, dberr := recv.Repo.SavePicture(picData, srcDomain)
	if dberr != nil {
		log.Warn().Err(err).Str("url", indieweb.Author.Picture).Msg("Unable to save downloaded author picture. Reverting to anonymous.")
		indieweb.Author.Picture = fmt.Sprintf("/pictures/%s", pictures.Anonymous)
		return
	}

	indieweb.Author.Picture = fmt.Sprintf("/pictures/%s", srcDomain)
}

func nonIndiewebTitle(body string, wm mf.Mention) string {
	titleMatch := titleRegexp.FindStringSubmatch(body)
	title := wm.Source
	if titleMatch != nil {
		title = titleMatch[1]
	}
	return title
}
