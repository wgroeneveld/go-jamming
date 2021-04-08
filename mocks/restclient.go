
package mocks

import (
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

