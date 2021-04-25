package send

import (
	"brainbaking.com/go-jamming/rest"
	"github.com/rs/zerolog/log"
	"net/url"
	"regexp"
	"strings"
	"willnorris.com/go/microformats"
)

const (
	typeWebmention string = "webmention"
	typeUnknown    string = "unknown"
	typePingback   string = "pingback"
)

var (
	relWebmention = regexp.MustCompile(`rel="??'??webmention`)
)

func (sndr *Sender) discover(target string) (link string, mentionType string) {
	mentionType = typeUnknown
	header, body, err := sndr.RestClient.GetBody(target)
	if err != nil {
		log.Warn().Str("target", target).Msg("Failed to discover possible endpoint, aborting send")
		return
	}
	link = header.Get(rest.RequestUrl) // default to a possible redirect of the target
	baseUrl, _ := url.Parse(link)

	// prefer links in the header over the html itself.
	for _, possibleLink := range header.Values("link") {
		if relWebmention.MatchString(possibleLink) {
			return buildWebmentionHeaderLink(possibleLink, baseUrl), typeWebmention
		}
	}
	if header.Get("X-Pingback") != "" {
		return header.Get("X-Pingback"), typePingback
	}

	// this also complies with w3.org regulations: relative endpoint could be possible
	format := microformats.Parse(strings.NewReader(body), baseUrl)
	if len(format.Rels[typeWebmention]) > 0 {
		mentionType = typeWebmention
		for _, possibleWm := range format.Rels[typeWebmention] {
			if possibleWm != link {
				link = possibleWm
				return
			}
		}
	} else if len(format.Rels[typePingback]) > 0 {
		mentionType = typePingback
		link = format.Rels[typePingback][0]
	}

	return
}

// buildWebmentionHeaderLink tries to extract the link from the link header.
// e.g. Link: <http://aaronpk.example/webmention-endpoint>; rel="webmention"
// could also be comma-separated, e.g. <https://webmention.rocks/test/19/webmention/error>; rel="other", <https://webmention.rocks/test/19/webmention?head=true>; rel="webmention"
func buildWebmentionHeaderLink(link string, baseUrl *url.URL) (wm string) {
	if strings.Contains(link, ",") {
		for _, possibleLink := range strings.Split(link, ",") {
			if relWebmention.MatchString(possibleLink) {
				link = strings.TrimSpace(possibleLink)
			}
		}
	}
	raw := strings.Split(link, ";")[0][1:]
	wm = raw[:len(raw)-1]
	if !strings.HasPrefix(wm, "http") {
		abs, _ := baseUrl.Parse(wm)
		wm = abs.String()
	}

	return
}
