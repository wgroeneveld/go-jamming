package mocks

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// neat trick! https://medium.com/@matryer/meet-moq-easily-mock-interfaces-in-go-476444187d10
type RestClientMock struct {
	GetFunc     func(string) (*http.Response, error)
	GetBodyFunc func(string) (string, error)
	PostFunc    func(string, string, string) error
}

// although these are still requied to match the rest.Client interface.
func (m *RestClientMock) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}
func (m *RestClientMock) GetBody(url string) (string, error) {
	return m.GetBodyFunc(url)
}

func (m *RestClientMock) Post(url string, contentType string, body string) error {
	return m.PostFunc(url, contentType, body)
}

func RelPathGetBodyFunc(t *testing.T, relPath string) func(string) (string, error) {
	return func(url string) (string, error) {
		// url: https://brainbaking.com/something-something.html
		// want: ../../mocks/something-something.html
		mockfile := relPath + strings.ReplaceAll(url, "https://brainbaking.com/", "")
		html, err := ioutil.ReadFile(mockfile)
		if err != nil {
			t.Error(err)
		}
		return string(html), nil
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
