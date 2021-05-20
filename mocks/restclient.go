package mocks

import (
	"brainbaking.com/go-jamming/rest"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// neat trick! https://medium.com/@matryer/meet-moq-easily-mock-interfaces-in-go-476444187d10
type RestClientMock struct {
	HeadFunc     func(string) (*http.Response, error)
	GetFunc      func(string) (*http.Response, error)
	GetBodyFunc  func(string) (http.Header, string, error)
	PostFunc     func(string, string, string) error
	PostFormFunc func(string, url.Values) error
}

// although these are still required to match the rest.Client interface.
func (m *RestClientMock) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}
func (m *RestClientMock) Head(url string) (*http.Response, error) {
	return m.HeadFunc(url)
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
		if key == "link1" || key == "link2" {
			key = "link"
		}
		httpHeader.Add(key, value.(string))
	}
	return httpHeader
}

func Head200ContentXml() func(string) (*http.Response, error) {
	return func(s string) (*http.Response, error) {
		return &http.Response{
			Header: map[string][]string{
				"Content-Type": {"text/xml"},
			},
			StatusCode: 200,
		}, nil
	}
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
			header := http.Header{}
			header.Set(rest.RequestUrl, url) // mimic actual implementation to track possible redirects
			return header, string(html), nil
		}
		headerJson := map[string]interface{}{}
		json.Unmarshal(headerData, &headerJson)

		header := toHttpHeader(headerJson)
		header.Set(rest.RequestUrl, url) // mimic actual implementation to track possible redirects
		return header, string(html), nil
	}
}
