package mf

import (
	"brainbaking.com/go-jamming/rest"
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

// TargetDomain parses the target url to extract the domain as part of the allowed webmention targets.
// This is the same as conf.FetchDomain(wm.Target), only without config, and without error handling.
// Assumes http(s) protocol, which should have been validated by now.
func (wm Mention) TargetDomain() string {
	return rest.Domain(wm.Target)
}

// SoureceDomain converts the Source to a domain name to be used in whitelisting/blacklisting (See TargetDomain()).
func (wm Mention) SourceDomain() string {
	return rest.Domain(wm.Source)
}

// Key returns a unique string representation of the mention for use in storage.
// TODO Profiling indicated that md5() consumes a lot of CPU power, so this could be replaced with db migration.
func (wm Mention) Key() string {
	return fmt.Sprintf("%x", md5.Sum([]byte("source="+wm.Source+",target="+wm.Target)))
}

func (wm Mention) SourceUrl() *url.URL {
	url, _ := url.Parse(wm.Source)
	return url
}
