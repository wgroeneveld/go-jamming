
package webmention

import (
	"fmt"
	"os"
	"crypto/md5"

	"github.com/wgroeneveld/go-jamming/common"
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

func (wm *webmention) asPath(conf *common.Config) string {
	filename := fmt.Sprintf("%x", md5.Sum([]byte("source=" + wm.source + ",target=" + wm.target)))
	domain, _ := conf.FetchDomain(wm.target)
	return conf.DataPath + "/" + domain + "/" + filename + ".json"
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

func (recv *receiver) processSourceBody(body string, wm webmention) {
	
}
