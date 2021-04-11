package send

import (
	"brainbaking.com/go-jamming/rest"
	"github.com/rs/zerolog/log"
	"strings"
	"willnorris.com/go/microformats"
)

const (
	TypeWebmention string = "webmention"
	TypeUnknown    string = "unknown"
	TypePingback   string = "pingback"
)

func (sndr *Sender) discover(target string) (link string, mentionType string) {
	mentionType = TypeUnknown
	header, body, err := sndr.RestClient.GetBody(target)
	if err != nil {
		log.Warn().Str("target", target).Msg("Failed to discover possible endpoint, aborting send")
		return
	}

	if strings.Contains(header.Get("link"), TypeWebmention) {
		return buildWebmentionHeaderLink(header.Get("link")), TypeWebmention
	}
	if header.Get("X-Pingback") != "" {
		return header.Get("X-Pingback"), TypePingback
	}

	// this also complies with w3.org regulations: relative endpoint could be possible
	format := microformats.Parse(strings.NewReader(body), rest.BaseUrlOf(target))
	if len(format.Rels[TypeWebmention]) > 0 {
		mentionType = TypeWebmention
		link = format.Rels[TypeWebmention][0]
	} else if len(format.Rels[TypePingback]) > 0 {
		mentionType = TypePingback
		link = format.Rels[TypePingback][0]
	}

	return
}

func buildWebmentionHeaderLink(link string) string {
	// e.g. Link: <http://aaronpk.example/webmention-endpoint>; rel="webmention"
	raw := strings.Split(link, ";")[0][1:]
	return raw[:len(raw)-1]
}
