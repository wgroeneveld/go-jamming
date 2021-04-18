package mf

import (
	"brainbaking.com/go-jamming/common"
	"crypto/md5"
	"fmt"
	"net/url"
)

// this should be passed along as a value object, not as a pointer
type Mention struct {
	Source string
	Target string
}

func (wm Mention) AsFormValues() url.Values {
	values := url.Values{}
	values.Add("source", wm.Source)
	values.Add("target", wm.Target)
	return values
}

func (wm Mention) String() string {
	return fmt.Sprintf("source: %s, target: %s", wm.Source, wm.Target)
}

func (wm Mention) Domain(conf *common.Config) string {
	domain, _ := conf.FetchDomain(wm.Target)
	return domain
}

func (wm Mention) Key() string {
	return fmt.Sprintf("%x", md5.Sum([]byte("source="+wm.Source+",target="+wm.Target)))
}

func (wm Mention) SourceUrl() *url.URL {
	url, _ := url.Parse(wm.Source)
	return url
}
