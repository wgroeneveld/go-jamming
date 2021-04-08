
package webmention

import (
	"fmt"
	"net/url"
	"strings"
	"os"
	"crypto/md5"

	"github.com/wgroeneveld/go-jamming/common"
	"github.com/wgroeneveld/go-jamming/rest"

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

type indiewebAuthor struct {
	name string
	picture string
}

type indiewebData struct {
	author indiewebAuthor
	name string
	content string
	published string // TODO to a date
	url string
	dateType string // TODO json property "type"
	source string
	target string
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
		indieweb = parseBodyAsNonIndiewebSite(body, wm)
	} else {
		indieweb = parseBodyAsIndiewebSite(hEntry, wm)
	}
	
	saveWebmentionToDisk(wm, indieweb)
	log.Info().Str("file", wm.asPath(recv.conf)).Msg("OK: webmention processed.")
}

func saveWebmentionToDisk(wm webmention, indieweb *indiewebData) {

}

// TODO I'm smelling very unstable code, apply https://golang.org/doc/effective_go#recover here?
// see https://github.com/willnorris/microformats/blob/main/microformats.go
func parseBodyAsIndiewebSite(hEntry *microformats.Microformat, wm webmention) *indiewebData {
	name := mfstr(hEntry, "name")
	authorName := mfstr(mfprop(hEntry, "author"), "name")
	if authorName == "" {
		authorName = mfprop(hEntry, "author").Value
	}
	// TODO sometimes it's picture.value??
	pic := mfstr(mfprop(hEntry, "author"), "photo")
	summary := mfstr(hEntry, "summary")
	contentEntry := mfmap(hEntry, "content")["value"]
	bridgyTwitterContent := mfstr(hEntry, "bridgy-twitter-content")

	return &indiewebData{
		name: name,
		author: indiewebAuthor{
			name: authorName,
			picture: pic,
		},
		content: determineContent(summary, contentEntry, bridgyTwitterContent),
		source: wm.source,
		target: wm.target,
	}

	//len(entry.Properties["hoopw"])
}

func shorten(txt string) string {
	if len(txt) <= 250 {
		return txt
	}
	return txt[0:250] + "..."
}

func determineContent(summary string, contentEntry string, bridgyTwitterContent string) string {
	if bridgyTwitterContent != "" {
		return shorten(bridgyTwitterContent)
	} else if summary != "" {
		return shorten(summary)
	}
	return shorten(contentEntry)
}

func parseBodyAsNonIndiewebSite(body string, wm webmention) *indiewebData {
	return nil
}
