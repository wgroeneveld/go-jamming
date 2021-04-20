package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// mimicing NotFound: https://golang.org/src/net/http/server.go?s=64787:64830#L2076
func BadRequest(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func TooManyRequests(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
}

func Unauthorized(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// Domain parses the target url to extract the domain as part of the allowed webmention targets.
// This is the same as conf.FetchDomain(wm.Target), only without config, and without error handling.
// Assumes http(s) protocol, which should have been validated by now.
func Domain(target string) string {
	withPossibleSubdomain := strings.Split(target, "/")[2]
	split := strings.Split(withPossibleSubdomain, ".")
	if len(split) == 2 {
		return withPossibleSubdomain // that was the extension, not the subdomain.
	}
	return fmt.Sprintf("%s.%s", split[1], split[2])
}

func Json(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(200)
	bytes, _ := json.MarshalIndent(data, "", "  ")
	w.Write(bytes)
}

func Accept(w http.ResponseWriter) {
	w.WriteHeader(202)
	w.Write([]byte("Thanks, bro. Will process this soon, pinky swear!"))
}

// assumes the URL is well-formed.
func BaseUrlOf(link string) *url.URL {
	obj, _ := url.Parse(link)
	baseUrl, _ := url.Parse(obj.Scheme + "://" + obj.Host)
	return baseUrl
}
