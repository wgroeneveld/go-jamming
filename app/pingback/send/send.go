package send

import (
	"github.com/rs/zerolog/log"
	"github.com/wgroeneveld/go-jamming/app/mf"
	"github.com/wgroeneveld/go-jamming/rest"
	"strings"
)

type xml string

var body xml = `<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>{$source}</string></value>
		</param>
		<param>
			<value><string>{$target}</string></value>
		</param>
	</params>
</methodCall>`

func (theXml xml) replace(key string, value string) xml {
	return xml(strings.ReplaceAll(theXml.String(), key, value))
}

func (theXml xml) String() string {
	return string(theXml)
}

func (theXml xml) fill(mention mf.Mention) string {
	return theXml.
		replace("{$source}", mention.Source).
		replace("{$target}", mention.Target).
		String()
}

type Sender struct {
	RestClient rest.Client
}

func (sender *Sender) SendPingbackToEndpoint(endpoint string, mention mf.Mention) {
	err := sender.RestClient.Post(endpoint, "text/xml", body.fill(mention))
	if err != nil {
		log.Err(err).Str("wm", mention.String()).Msg("Unable to send pingback")
	}
	log.Info().Str("wm", mention.String()).Msg("Pingback sent")
}
