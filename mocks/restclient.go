
package mocks

import (
	"testing"
	"io/ioutil"
	"net/http"
)

// neat trick! https://medium.com/@matryer/meet-moq-easily-mock-interfaces-in-go-476444187d10
type RestClientMock struct {
	GetFunc func(string) (*http.Response, error)
	GetBodyFunc func(string) (string, error)
}

// although these are still requied to match the rest.Client interface. 
func (m *RestClientMock) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}
func (m *RestClientMock) GetBody(url string) (string, error) {
	return m.GetBodyFunc(url)
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
