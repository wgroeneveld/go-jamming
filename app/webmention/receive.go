
package webmention

import (
	"fmt"

	"github.com/wgroeneveld/go-jamming/common"

	"github.com/rs/zerolog/log"
)

type webmention struct {
	source string
	target string
}

func (wm *webmention) String() string {
    return fmt.Sprintf("source: %s, target: %s", wm.source, wm.target)
}

func (wm *webmention) receive() {
	log.Info().Str("webmention", wm.String()).Msg("OK: looks valid")
	body, geterr := common.Get(wm.source)

	if geterr != nil {
        log.Warn().Str("source", wm.source).Msg("  ABORT: invalid url")
		wm.deletePossibleOlderWebmention()
		return
	}

	wm.processSourceBody(body)	
}

func (wm *webmention) deletePossibleOlderWebmention() {

}

func (wm *webmention) processSourceBody(body string) {
	
}
