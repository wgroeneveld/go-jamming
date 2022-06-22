package recv

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/notifier"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"willnorris.com/go/microformats"
)

type Receiver struct {
	RestClient rest.Client
	Notifier   notifier.Notifier
	Conf       *common.Config
	Repo       db.MentionRepo
}

var (
	titleRegexp = regexp.MustCompile(`<title>(.*?)<\/title>`)

	errPicUnableToDownload          = errors.New("Unable to download author picture")
	errPicNoRealImage               = errors.New("Downloaded author picture is not a real image")
	errPicUnableToSave              = errors.New("Unable to save downloaded author picture")
	errWontDownloadBecauseOfPrivacy = errors.New("Will not save locally because it's from a silo domain")
)

func (recv *Receiver) Receive(wm mf.Mention) {
	if recv.Conf.IsBlacklisted(wm.Source) {
		log.Warn().Stringer("wm", wm).Msg("  ABORT: source url comes from blacklisted domain!")
		return
	}

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
	defer func() {
		if r := recover(); r != nil {
			log.Error().Str("panic", fmt.Sprintf("%q", r)).Stringer("wm", wm).Msg("ABORT: panic recovery while processing wm")
		}
	}()

	data := microformats.Parse(strings.NewReader(body), wm.SourceUrl())
	indieweb := recv.convertBodyToIndiewebData(body, wm, data)
	recv.ProcessAuthorPicture(indieweb)

	if recv.Conf.IsWhitelisted(wm.Source) {
		recv.ProcessWhitelistedMention(wm, indieweb)
	} else {
		recv.ProcessMentionInModeration(wm, indieweb)
	}
}

func (recv *Receiver) ProcessMentionInModeration(wm mf.Mention, indieweb *mf.IndiewebData) {
	key, err := recv.Repo.InModeration(wm, indieweb)
	if err != nil {
		log.Error().Err(err).Stringer("wm", wm).Msg("Failed to save new mention to in moderation db")
	}
	err = recv.Notifier.NotifyInModeration(wm, indieweb)
	if err != nil {
		log.Error().Err(err).Msg("Failed to notify")
	}
	log.Info().Str("key", key).Msg("OK: Webmention processed, in moderation.")
}

func (recv *Receiver) ProcessWhitelistedMention(wm mf.Mention, indieweb *mf.IndiewebData) {
	key, err := recv.Repo.Save(wm, indieweb)
	if err != nil {
		log.Error().Err(err).Stringer("wm", wm).Msg("Failed to save new mention to db")
	}
	err = recv.Notifier.NotifyReceived(wm, indieweb)
	if err != nil {
		log.Error().Err(err).Msg("Failed to notify")
	}
	log.Info().Str("key", key).Msg("OK: Webmention processed, in whitelist.")
}

func (recv *Receiver) ProcessAuthorPicture(indieweb *mf.IndiewebData) {
	if indieweb.Author.Picture != "" {
		err := recv.saveAuthorPictureLocally(indieweb)
		if err != nil {
			log.Error().Err(err).Str("url", indieweb.Author.Picture).Msg("Failed to save picture. Reverting to anonymous")
			indieweb.Author.AnonymizePicture()

			if err == errWontDownloadBecauseOfPrivacy {
				indieweb.Author.AnonymizeName()
			}
		}
	} else {
		indieweb.Author.AnonymizePicture()
	}
}

func (recv *Receiver) convertBodyToIndiewebData(body string, wm mf.Mention, mfRoot *microformats.Data) *mf.IndiewebData {
	hEntry := mf.HEntry(mfRoot)
	hCard := mf.HCard(mfRoot)
	if hEntry == nil {
		return recv.parseBodyAsNonIndiewebSite(body, wm)
	}
	return recv.parseBodyAsIndiewebSite(hEntry, hCard, wm)
}

// see https://github.com/willnorris/microformats/blob/main/microformats.go
func (recv *Receiver) parseBodyAsIndiewebSite(hEntry *microformats.Microformat, hCard *microformats.Microformat, wm mf.Mention) *mf.IndiewebData {
	return &mf.IndiewebData{
		Name:         mf.Str(hEntry, "name"),
		Author:       mf.NewAuthor(hEntry, hCard),
		Content:      mf.Content(hEntry),
		Url:          mf.Url(hEntry, wm.Source),
		Published:    mf.Published(hEntry),
		Source:       wm.Source,
		Target:       wm.Target,
		IndiewebType: mf.Type(hEntry),
	}
}

func (recv *Receiver) parseBodyAsNonIndiewebSite(body string, wm mf.Mention) *mf.IndiewebData {
	title := nonIndiewebTitle(body, wm)
	return &mf.IndiewebData{
		Author: mf.IndiewebAuthor{
			Name: wm.Source,
		},
		Name:         title,
		Content:      title,
		Published:    mf.PublishedNow(),
		Url:          wm.Source,
		IndiewebType: mf.TypeMention,
		Source:       wm.Source,
		Target:       wm.Target,
	}
}

// saveAuthorPictureLocally tries to download the author picture and checks if it's valid based on img header.
// If it succeeds, it alters the picture path to a local /pictures/x one.
// If it fails, it returns an error.
// If strict is true, this refuses to download from silo sources such as brid.gy because of privacy concerns.
func (recv *Receiver) saveAuthorPictureLocally(indieweb *mf.IndiewebData) error {
	srcDomain := rest.Domain(indieweb.Source)
	if common.Includes(rest.SiloDomains, srcDomain) {
		return errWontDownloadBecauseOfPrivacy
	}

	_, picData, err := recv.RestClient.GetBody(indieweb.Author.Picture)
	if err != nil {
		return errPicUnableToDownload
	}
	if len(picData) < 8 || !rest.IsRealImage([]byte(picData[0:8])) {
		return errPicNoRealImage
	}

	_, dberr := recv.Repo.SavePicture(picData, srcDomain)
	if dberr != nil {
		return errPicUnableToSave
	}

	indieweb.Author.Picture = fmt.Sprintf("/pictures/%s", srcDomain)
	return nil
}

func nonIndiewebTitle(body string, wm mf.Mention) string {
	titleMatch := titleRegexp.FindStringSubmatch(body)
	title := wm.Source
	if titleMatch != nil {
		title = titleMatch[1]
	}
	return title
}
