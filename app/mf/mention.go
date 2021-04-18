package mf

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strings"
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

// Domain parses the target url to extract the domain as part of the allowed webmention targets.
// This is the same as conf.FetchDomain(wm.Target), only without config, and without error handling.
// Assumes http(s) protocol, which should have been validated by now.
func (wm Mention) Domain() string {
	withPossibleSubdomain := strings.Split(wm.Target, "/")[2]
	split := strings.Split(withPossibleSubdomain, ".")
	if len(split) == 2 {
		return withPossibleSubdomain // that was the extention, not the subdomain.
	}
	return fmt.Sprintf("%s.%s", split[1], split[2])
}

func (wm Mention) Key() string {
	return fmt.Sprintf("%x", md5.Sum([]byte("source="+wm.Source+",target="+wm.Target)))
}

func (wm Mention) SourceUrl() *url.URL {
	url, _ := url.Parse(wm.Source)
	return url
}
