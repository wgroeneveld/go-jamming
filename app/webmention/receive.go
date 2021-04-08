
package webmention

import (
	"fmt"

	"github.com/wgroeneveld/go-jamming/rest"

	"github.com/rs/zerolog/log"
)

type webmention struct {
	source string
	target string
}

func (wm *webmention) String() string {
    return fmt.Sprintf("source: %s, target: %s", wm.source, wm.target)
}

// used as a "class" to iject dependencies, just to be able to test. Do NOT like htis. 
// Is there a better way? e.g. in validate, I just pass rest.Client as an arg. Not great either. 
type receiver struct {
	restClient rest.Client
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

}

func (recv *receiver) processSourceBody(body string, wm webmention) {
	
}
