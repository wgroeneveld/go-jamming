package send

import (
	"brainbaking.com/go-jamming/rest"
	"github.com/rs/zerolog/log"
	"strings"
	"willnorris.com/go/microformats"
)

const (
	typeWebmention string = "webmention"
	typeUnknown    string = "unknown"
	typePingback   string = "pingback"
)

func (sndr *Sender) discover(target string) (link string, mentionType string) {
	mentionType = typeUnknown
	header, body, err := sndr.RestClient.GetBody(target)
	if err != nil {
		log.Warn().Str("target", target).Msg("Failed to discover possible endpoint, aborting send")
		return
	}

	if strings.Contains(header.Get("link"), typeWebmention) {
		return buildWebmentionHeaderLink(header.Get("link")), typeWebmention
	}
	if header.Get("X-Pingback") != "" {
		return header.Get("X-Pingback"), typePingback
	}

	// this also complies with w3.org regulations: relative endpoint could be possible
	format := microformats.Parse(strings.NewReader(body), rest.BaseUrlOf(target))
	if len(format.Rels[typeWebmention]) > 0 {
		mentionType = typeWebmention
		link = format.Rels[typeWebmention][0]
	} else if len(format.Rels[typePingback]) > 0 {
		mentionType = typePingback
		link = format.Rels[typePingback][0]
	}

	return
}

// e.g. Link: <http://aaronpk.example/webmention-endpoint>; rel="webmention"
func buildWebmentionHeaderLink(link string) string {
	raw := strings.Split(link, ";")[0][1:]
	return raw[:len(raw)-1]
}
