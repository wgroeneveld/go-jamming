
package webmention

import (
	"strings"

	"github.com/wgroeneveld/go-jamming/common"
)

func isValidUrl(url string) bool {
	return url != "" &&
		(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://"))
}

func isValidDomain(url string, conf *common.Config) bool {
	for _, domain := range conf.AllowedWebmentionSources {
		if strings.Index(url, domain) != -1 {
			return true
		}
	}
	return false
}

type httpReq interface {
	FormValue(key string) string
}
type httpHeader interface {
	Get(key string) string
}

func validate(r httpReq, h httpHeader, conf *common.Config) bool {
	return h.Get("Content-Type") == "application/x-www-form-urlencoded" &&
		isValidUrl(r.FormValue("source")) &&
		isValidUrl(r.FormValue("target")) &&
		r.FormValue("source") != r.FormValue("target") &&
		isValidDomain(r.FormValue("target"), conf)
}
