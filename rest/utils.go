package rest

import (
	"net/http"
	"net/url"
)

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "400 bad request", http.StatusBadRequest)
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
