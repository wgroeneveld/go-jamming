package recv

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"willnorris.com/go/microformats"
)

// used as a "class" to iject dependencies, just to be able to test. Do NOT like htis.
// Is there a better way? e.g. in validate, I just pass rest.Client as an arg. Not great either.
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

func nonIndiewebTitle(body string, wm mf.Mention) string {
	titleMatch := titleRegexp.FindStringSubmatch(body)
	title := wm.Source
	if titleMatch != nil {
		title = titleMatch[1]
	}
	return title
}
