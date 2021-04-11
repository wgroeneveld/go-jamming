package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	Get(url string) (*http.Response, error)
	Post(url string, contentType string, body string) error
	GetBody(url string) (http.Header, string, error)
	PostForm(url string, formData url.Values) error
}

type HttpClient struct {
}

func (client *HttpClient) PostForm(url string, formData url.Values) error {
	resp, err := http.PostForm(url, formData)
	if err != nil {
		return err
	}
	if !isStatusOk(resp) {
		return fmt.Errorf("POST Form to %s: Status code is not OK (%d)", url, resp.StatusCode)
	}
	return nil
}

func (client *HttpClient) Post(url string, contenType string, body string) error {
	resp, err := http.Post(url, contenType, strings.NewReader(body))
	if err != nil {
		return err
	}
	if !isStatusOk(resp) {
		return fmt.Errorf("POST to %s: Status code is not OK (%d)", url, resp.StatusCode)
	}
	return nil
}

// something like this? https://freshman.tech/snippets/go/http-response-to-string/
func (client *HttpClient) GetBody(url string) (http.Header, string, error) {
	resp, geterr := client.Get(url)
	if geterr != nil {
		return nil, "", geterr
	}

	body, err := ReadBodyFromResponse(resp)
	if err != nil {
		return nil, "", err
	}

	return resp.Header, body, nil
}

func ReadBodyFromResponse(resp *http.Response) (string, error) {
	if !isStatusOk(resp) {
		return "", fmt.Errorf("Status code is not OK (%d)", resp.StatusCode)
	}

	body, readerr := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if readerr != nil {
		return "", readerr
	}
	return string(body), nil
}

func isStatusOk(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode <= 299
}

func (client *HttpClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}
