
package webmention

import (
	"strings"

	"github.com/wgroeneveld/go-jamming/common"
	"github.com/wgroeneveld/go-jamming/rest"

	"github.com/rs/zerolog/log"
)

func isValidUrl(url string) bool {
	return url != "" &&
		(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://"))
}

func isValidDomain(url string, conf *common.Config) bool {
	_, err := conf.FetchDomain(url)
	if err != nil {
		return false
	}
	return true
}

func isValidTargetUrl(url string, httpClient rest.Client) bool {
	_, err := httpClient.Get(url)
	if err != nil {
		log.Warn().Str("target", url).Msg("Invalid target URL")
		return false
	}
	return true
}

func validate(r rest.HttpReq, h rest.HttpHeader, conf *common.Config) bool {
	return h.Get("Content-Type") == "application/x-www-form-urlencoded" &&
		isValidUrl(r.FormValue("source")) &&
		isValidUrl(r.FormValue("target")) &&
		r.FormValue("source") != r.FormValue("target") &&
		isValidDomain(r.FormValue("target"), conf)
}
