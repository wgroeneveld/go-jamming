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

type imageType []byte

var (
	jpg                 = imageType{0xFF, 0xD8}
	bmp                 = imageType{0x42, 0x4D}
	gif                 = imageType{0x47, 0x49, 0x46}
	png                 = imageType{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	tiffI               = imageType{0x49, 0x49, 0x2A, 0x00}
	tiffM               = imageType{0x4D, 0x4D, 0x00, 0x2A}
	webp                = imageType{0x52, 0x49, 0x46, 0x46} // RIFF 32 bits
	supportedImageTypes = []imageType{jpg, png, gif, bmp, webp, tiffI, tiffM}

	// SiloDomains are domains where mentions of multiple individuals may come from.
	// These are privacy issues and will be anonymized as such.
	SiloDomains = []string{"brid.gy", "twitter.com", "facebook.com"}
)

// IsRealImage checks the first few bytes of the provided data to see if it's a real image.
// Image headers supported: gif/jpg/png/webp/bmp
func IsRealImage(data []byte) bool {
	if len(data) < 8 {
		return false
	}

	for _, imgType := range supportedImageTypes {
		checkedBits := 0
		for i, bit := range imgType {
			if data[i] == bit {
				checkedBits++
			}
		}
		if checkedBits == len(imgType) {
			return true
		}
	}
	return false
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
