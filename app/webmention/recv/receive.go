package recv

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
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
}

func (recv *Receiver) Receive(wm mf.Mention) {
	log.Info().Str("Webmention", wm.String()).Msg("OK: looks valid")
	body, geterr := recv.RestClient.GetBody(wm.Source)

	if geterr != nil {
		log.Warn().Str("source", wm.Source).Msg("  ABORT: invalid url")
		recv.deletePossibleOlderWebmention(wm)
		return
	}

	recv.processSourceBody(body, wm)
}

func (recv *Receiver) deletePossibleOlderWebmention(wm mf.Mention) {
	os.Remove(wm.AsPath(recv.Conf))
}

func getHEntry(data *microformats.Data) *microformats.Microformat {
	for _, itm := range data.Items {
		if common.Includes(itm.Type, "h-entry") {
			return itm
		}
	}
	return nil
}

func (recv *Receiver) processSourceBody(body string, wm mf.Mention) {
	if !strings.Contains(body, wm.Target) {
		log.Warn().Str("target", wm.Target).Msg("ABORT: no mention of target found in html src of source!")
		return
	}

	data := microformats.Parse(strings.NewReader(body), wm.SourceUrl())
	indieweb := recv.convertBodyToIndiewebData(body, wm, getHEntry(data))

	recv.saveWebmentionToDisk(wm, indieweb)
	log.Info().Str("file", wm.AsPath(recv.Conf)).Msg("OK: Webmention processed.")
}

func (recv *Receiver) convertBodyToIndiewebData(body string, wm mf.Mention, hEntry *microformats.Microformat) *mf.IndiewebData {
	if hEntry == nil {
		return recv.parseBodyAsNonIndiewebSite(body, wm)
	}
	return recv.parseBodyAsIndiewebSite(hEntry, wm)
}

func (recv *Receiver) saveWebmentionToDisk(wm mf.Mention, indieweb *mf.IndiewebData) {
	jsonData, jsonErr := json.Marshal(indieweb)
	if jsonErr != nil {
		log.Err(jsonErr).Msg("Unable to serialize Webmention into JSON")
	}
	err := ioutil.WriteFile(wm.AsPath(recv.Conf), jsonData, fs.ModePerm)
	if err != nil {
		log.Err(err).Msg("Unable to save Webmention to disk")
	}
}

// TODO I'm smelling very unstable code, apply https://golang.org/doc/effective_go#recover here?
// see https://github.com/willnorris/microformats/blob/main/microformats.go
func (recv *Receiver) parseBodyAsIndiewebSite(hEntry *microformats.Microformat, wm mf.Mention) *mf.IndiewebData {
	name := mf.Str(hEntry, "name")
	pic := mf.Str(mf.Prop(hEntry, "author"), "photo")
	mfType := mf.DetermineType(hEntry)

	return &mf.IndiewebData{
		Name: name,
		Author: mf.IndiewebAuthor{
			Name:    mf.DetermineAuthorName(hEntry),
			Picture: pic,
		},
		Content:      mf.DetermineContent(hEntry),
		Url:          mf.DetermineUrl(hEntry, wm.Source),
		Published:    mf.DeterminePublishedDate(hEntry, recv.Conf.UtcOffset),
		Source:       wm.Source,
		Target:       wm.Target,
		IndiewebType: mfType,
	}
}

func (recv *Receiver) parseBodyAsNonIndiewebSite(body string, wm mf.Mention) *mf.IndiewebData {
	r := regexp.MustCompile(`<title>(.*?)<\/title>`)
	titleMatch := r.FindStringSubmatch(body)
	title := wm.Source
	if titleMatch != nil {
		title = titleMatch[1]
	}
	return &mf.IndiewebData{
		Author: mf.IndiewebAuthor{
			Name: wm.Source,
		},
		Name:         title,
		Content:      title,
		Published:    mf.PublishedNow(recv.Conf.UtcOffset),
		Url:          wm.Source,
		IndiewebType: "mention",
		Source:       wm.Source,
		Target:       wm.Target,
	}
}
