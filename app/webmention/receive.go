
package webmention

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/wgroeneveld/go-jamming/common"
	"github.com/wgroeneveld/go-jamming/rest"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"willnorris.com/go/microformats"
)

type Mention struct {
	Source string
	Target string
}

func (wm *Mention) String() string {
    return fmt.Sprintf("source: %s, target: %s", wm.Source, wm.Target)
}

func (wm *Mention) asPath(conf *common.Config) string {
	filename := fmt.Sprintf("%x", md5.Sum([]byte("source=" + wm.Source+ ",target=" + wm.Target)))
	domain, _ := conf.FetchDomain(wm.Target)
	return conf.DataPath + "/" + domain + "/" + filename + ".json"
}

func (wm *Mention) sourceUrl() *url.URL {
	url, _ := url.Parse(wm.Source)
	return url
}

// used as a "class" to iject dependencies, just to be able to test. Do NOT like htis. 
// Is there a better way? e.g. in validate, I just pass rest.Client as an arg. Not great either. 
type Receiver struct {
	RestClient rest.Client
	Conf       *common.Config
}

func (recv *Receiver) Receive(wm Mention) {
	log.Info().Str("Webmention", wm.String()).Msg("OK: looks valid")
	body, geterr := recv.RestClient.GetBody(wm.Source)

	if geterr != nil {
        log.Warn().Str("source", wm.Source).Msg("  ABORT: invalid url")
		recv.deletePossibleOlderWebmention(wm)
		return
	}

	recv.processSourceBody(body, wm)
}

func (recv *Receiver) deletePossibleOlderWebmention(wm Mention) {
	os.Remove(wm.asPath(recv.Conf))
}

func getHEntry(data *microformats.Data) *microformats.Microformat {
	for _, itm := range data.Items {
		if common.Includes(itm.Type, "h-entry") {
			return itm
		}
	}
	return nil
}


func (recv *Receiver) processSourceBody(body string, wm Mention) {
	if !strings.Contains(body, wm.Target) {
		log.Warn().Str("target", wm.Target).Msg("ABORT: no mention of target found in html src of source!")
		return
	}

	data := microformats.Parse(strings.NewReader(body), wm.sourceUrl())
	indieweb := recv.convertBodyToIndiewebData(body, wm, getHEntry(data))

	recv.saveWebmentionToDisk(wm, indieweb)
	log.Info().Str("file", wm.asPath(recv.Conf)).Msg("OK: Webmention processed.")
}

func (recv *Receiver) convertBodyToIndiewebData(body string, wm Mention, hEntry *microformats.Microformat) *indiewebData {
	if hEntry == nil {
		return recv.parseBodyAsNonIndiewebSite(body, wm)
	}
	return recv.parseBodyAsIndiewebSite(hEntry, wm)
}

func (recv *Receiver) saveWebmentionToDisk(wm Mention, indieweb *indiewebData) {
	jsonData, jsonErr := json.Marshal(indieweb)
	if jsonErr != nil {
		log.Err(jsonErr).Msg("Unable to serialize Webmention into JSON")
	}
	err := ioutil.WriteFile(wm.asPath(recv.Conf), jsonData, fs.ModePerm)
	if err != nil {
		log.Err(err).Msg("Unable to save Webmention to disk")
	}
}

// TODO I'm smelling very unstable code, apply https://golang.org/doc/effective_go#recover here?
// see https://github.com/willnorris/microformats/blob/main/microformats.go
func (recv *Receiver) parseBodyAsIndiewebSite(hEntry *microformats.Microformat, wm Mention) *indiewebData {
	name := mfStr(hEntry, "name")
	pic := mfStr(mfProp(hEntry, "author"), "photo")
	mfType := determineMfType(hEntry)

	return &indiewebData{
		Name: name,
		Author: indiewebAuthor{
			Name: determineAuthorName(hEntry),
			Picture: pic,
		},
		Content: determineContent(hEntry),
		Url: determineUrl(hEntry, wm.Source),
		Published: determinePublishedDate(hEntry, recv.Conf.UtcOffset),
		Source: wm.Source,
		Target: wm.Target,
		IndiewebType: mfType,
	}
}

func (recv *Receiver) parseBodyAsNonIndiewebSite(body string, wm Mention) *indiewebData {
	r := regexp.MustCompile(`<title>(.*?)<\/title>`)
	titleMatch := r.FindStringSubmatch(body)
	title := wm.Source
	if titleMatch != nil {
		title = titleMatch[1]
	}
	return &indiewebData{
		Author: indiewebAuthor{
			Name: wm.Source,
		},
		Name: title,
		Content: title,
		Published: publishedNow(recv.Conf.UtcOffset),
		Url: wm.Source,
		IndiewebType: "mention",
		Source: wm.Source,
		Target: wm.Target,
	}
}
