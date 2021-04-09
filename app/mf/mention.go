package mf

import (
	"crypto/md5"
	"fmt"
	"github.com/wgroeneveld/go-jamming/common"
	"net/url"
)

type Mention struct {
	Source string
	Target string
}

func (wm *Mention) String() string {
	return fmt.Sprintf("source: %s, target: %s", wm.Source, wm.Target)
}

func (wm *Mention) AsPath(conf *common.Config) string {
	filename := fmt.Sprintf("%x", md5.Sum([]byte("source=" + wm.Source+ ",target=" + wm.Target)))
	domain, _ := conf.FetchDomain(wm.Target)
	return conf.DataPath + "/" + domain + "/" + filename + ".json"
}

func (wm *Mention) SourceUrl() *url.URL {
	url, _ := url.Parse(wm.Source)
	return url
}
