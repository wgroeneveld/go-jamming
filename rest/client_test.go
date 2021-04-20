package rest

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

var client = HttpClient{}

func TestGetBodyWithinLimitsReturnsHeadersAndBodyString(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("amaigat", "jup.")
		w.WriteHeader(200)
		data, err := ioutil.ReadFile("../mocks/samplerss.xml") // is about 1.6 MB
		assert.NoError(t, err)
		w.Write(data)
	})
	srv := &http.Server{Addr: ":6666", Handler: mux}
	defer srv.Close()

	go func() {
		srv.ListenAndServe()
	}()
	headers, body, err := client.GetBody("http://localhost:6666/")

	assert.NoError(t, err)
	assert.Equal(t, "jup.", headers.Get("amaigat"))
	assert.Contains(t, body, "<rss")
}

func TestGetBodyOf404ReturnsError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	srv := &http.Server{Addr: ":6666", Handler: mux}
	defer srv.Close()

	go func() {
		srv.ListenAndServe()
	}()
	_, body, err := client.GetBody("http://localhost:6666/")

	assert.Contains(t, err.Error(), "404")
	assert.Equal(t, "", body)
}

func TestGetBodyOfTooLargeContentReturnsError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		garbage := make([]byte, MaxBytes*2)
		for i := 0; i < len(garbage); i++ {
			garbage[i] = 'A'
		}
		w.Write(garbage)
	})
	srv := &http.Server{Addr: ":6666", Handler: mux}
	defer srv.Close()

	go func() {
		srv.ListenAndServe()
	}()
	_, body, err := client.GetBody("http://localhost:6666/")

	assert.Equal(t, ResponseAboveLimit, errors.Unwrap(err))
	assert.Equal(t, "", body)
}
