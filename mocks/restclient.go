package mocks

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// neat trick! https://medium.com/@matryer/meet-moq-easily-mock-interfaces-in-go-476444187d10
type RestClientMock struct {
	GetFunc      func(string) (*http.Response, error)
	GetBodyFunc  func(string) (http.Header, string, error)
	PostFunc     func(string, string, string) error
	PostFormFunc func(string, url.Values) error
}

// although these are still requied to match the rest.Client interface.
func (m *RestClientMock) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}
func (m *RestClientMock) GetBody(url string) (http.Header, string, error) {
	return m.GetBodyFunc(url)
}

func (client *RestClientMock) PostForm(url string, formData url.Values) error {
	return client.PostFormFunc(url, formData)
}

func (m *RestClientMock) Post(url string, contentType string, body string) error {
	return m.PostFunc(url, contentType, body)
}

func toHttpHeader(header map[string]interface{}) http.Header {
	httpHeader := http.Header{}
	for key, value := range header {
		httpHeader.Add(key, value.(string))
	}
	return httpHeader
}

func RelPathGetBodyFunc(relPath string) func(string) (http.Header, string, error) {
	return func(url string) (http.Header, string, error) {
		log.Debug().Str("url", url).Msg("  - GET call")
		// url: https://brainbaking.com/something-something.html
		// want: ../../mocks/something-something.html
		mockfile := relPath + strings.ReplaceAll(url, "https://brainbaking.com/", "")
		html, err := ioutil.ReadFile(mockfile)
		if err != nil {
			return nil, "", err
		}

		headerData, headerFileErr := ioutil.ReadFile(strings.ReplaceAll(mockfile, ".html", "-headers.json"))
		if headerFileErr != nil {
			return http.Header{}, string(html), nil
		}
		headerJson := map[string]interface{}{}
		json.Unmarshal(headerData, &headerJson)

		return toHttpHeader(headerJson), string(html), nil
	}
}

func BodyFunc(t *testing.T, mockfile string) func(string) (string, error) {
	html, err := ioutil.ReadFile(mockfile)
	if err != nil {
		t.Error(err)
	}
	return func(url string) (string, error) {
		return string(html), nil
	}
}
