
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

type webmention struct {
	source string
	target string
}

func (wm *webmention) String() string {
    return fmt.Sprintf("source: %s, target: %s", wm.source, wm.target)
}

func (wm *webmention) asPath(conf *common.Config) string {
	filename := fmt.Sprintf("%x", md5.Sum([]byte("source=" + wm.source + ",target=" + wm.target)))
	domain, _ := conf.FetchDomain(wm.target)
	return conf.DataPath + "/" + domain + "/" + filename + ".json"
}

func (wm *webmention) sourceUrl() *url.URL {
	url, _ := url.Parse(wm.source)
	return url
}

// used as a "class" to iject dependencies, just to be able to test. Do NOT like htis. 
// Is there a better way? e.g. in validate, I just pass rest.Client as an arg. Not great either. 
type receiver struct {
	restClient rest.Client
	conf *common.Config
}

func (recv *receiver) receive(wm webmention) {
	log.Info().Str("webmention", wm.String()).Msg("OK: looks valid")
	body, geterr := recv.restClient.GetBody(wm.source)

	if geterr != nil {
        log.Warn().Str("source", wm.source).Msg("  ABORT: invalid url")
		recv.deletePossibleOlderWebmention(wm)
		return
	}

	recv.processSourceBody(body, wm)
}

func (recv *receiver) deletePossibleOlderWebmention(wm webmention) {
	os.Remove(wm.asPath(recv.conf))
}

func getHEntry(data *microformats.Data) *microformats.Microformat {
	for _, itm := range data.Items {
		if common.Includes(itm.Type, "h-entry") {
			return itm
		}
	}
	return nil
}


func (recv *receiver) processSourceBody(body string, wm webmention) {
	if !strings.Contains(body, wm.target) {
		log.Warn().Str("target", wm.target).Msg("ABORT: no mention of target found in html src of source!")
		return
	}

	r := strings.NewReader(body)
	data := microformats.Parse(r, wm.sourceUrl())
	hEntry := getHEntry(data)
	var indieweb *indiewebData
	if hEntry == nil {
		indieweb = recv.parseBodyAsNonIndiewebSite(body, wm)
	} else {
		indieweb = recv.parseBodyAsIndiewebSite(hEntry, wm)
	}
	
	recv.saveWebmentionToDisk(wm, indieweb)
	log.Info().Str("file", wm.asPath(recv.conf)).Msg("OK: webmention processed.")
}

func (recv *receiver) saveWebmentionToDisk(wm webmention, indieweb *indiewebData) {
	jsonData, jsonErr := json.Marshal(indieweb)
	if jsonErr != nil {
		log.Err(jsonErr).Msg("Unable to serialize webmention into JSON")
	}
	err := ioutil.WriteFile(wm.asPath(recv.conf), jsonData, fs.ModePerm)
	if err != nil {
		log.Err(err).Msg("Unable to save webmention to disk")
	}
}

// TODO I'm smelling very unstable code, apply https://golang.org/doc/effective_go#recover here?
// see https://github.com/willnorris/microformats/blob/main/microformats.go
func (recv *receiver) parseBodyAsIndiewebSite(hEntry *microformats.Microformat, wm webmention) *indiewebData {
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
		Url: determineUrl(hEntry, wm.source),
		Published: determinePublishedDate(hEntry, recv.conf.UtcOffset),
		Source: wm.source,
		Target: wm.target,
		IndiewebType: mfType,
	}
}

func determinePublishedDate(hEntry *microformats.Microformat, utcOffset int) string {
	publishedDate := mfStr(hEntry, "published")
	if publishedDate == "" {
		return publishedNow(utcOffset)
	}
	return publishedDate
}

func determineAuthorName(hEntry *microformats.Microformat) string {
	authorName := mfStr(mfProp(hEntry, "author"), "name")
	if authorName == "" {
		return mfProp(hEntry, "author").Value
	}
	return authorName
}

func determineMfType(hEntry *microformats.Microformat) string {
	likeOf := mfStr(hEntry, "like-of")
	if likeOf != "" {
		return "like"
	}
	bookmarkOf := mfStr(hEntry, "bookmark-of")
	if bookmarkOf != "" {
		return "bookmark"
	}
	return "mention"
}

// Mastodon uids start with "tag:server", but we do want indieweb uids from other sources
func determineUrl(hEntry *microformats.Microformat, source string) string {
	uid := mfStr(hEntry, "uid")
	if uid != "" && strings.HasPrefix(uid, "http") {
		return uid
	}
	url := mfStr(hEntry, "url")
	if url != "" {
		return url
	}
	return source
}

func determineContent(hEntry *microformats.Microformat) string {
	bridgyTwitterContent := mfStr(hEntry, "bridgy-twitter-content")
	if bridgyTwitterContent != "" {
		return shorten(bridgyTwitterContent)
	}
	summary := mfStr(hEntry, "summary")
	if summary != "" {
		return shorten(summary)
	}
	contentEntry := mfMap(hEntry, "content")["value"]
	return shorten(contentEntry)
}

func (recv *receiver) parseBodyAsNonIndiewebSite(body string, wm webmention) *indiewebData {
	r := regexp.MustCompile(`<title>(.*?)<\/title>`)
	titleMatch := r.FindStringSubmatch(body)
	title := wm.source
	if titleMatch != nil {
		title = titleMatch[1]
	}
	return &indiewebData{
		Author: indiewebAuthor{
			Name: wm.source,
		},
		Name: title,
		Content: title,
		Published: publishedNow(recv.conf.UtcOffset),
		Url: wm.source,
		IndiewebType: "mention",
		Source: wm.source,
		Target: wm.target,
	}
}
