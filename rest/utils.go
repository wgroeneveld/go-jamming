package rest

import (
	"encoding/json"
	"net/http"
	"net/url"
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
