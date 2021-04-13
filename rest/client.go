package rest

import (
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client interface {
	Get(url string) (*http.Response, error)
	Post(url string, contentType string, body string) error
	GetBody(url string) (http.Header, string, error)
	PostForm(url string, formData url.Values) error
}

type HttpClient struct {
}

var (
	// do not use retryablehttp default impl - inject own logger and retry policies
	jammingHttp = &retryablehttp.Client{
		HTTPClient:   cleanhttp.DefaultPooledClient(),
		Logger:       &zeroLogWrapper{},
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
		RetryMax:     5,
		CheckRetry:   retryablehttp.DefaultRetryPolicy,
		Backoff:      retryablehttp.DefaultBackoff,
	}
)

func (client *HttpClient) PostForm(url string, formData url.Values) error {
	resp, err := jammingHttp.PostForm(url, formData)
	if err != nil {
		return fmt.Errorf("POST Form to %s: %v", url, err)
	}
	if !isStatusOk(resp) {
		return fmt.Errorf("POST Form to %s: Status code is not OK (%d)", url, resp.StatusCode)
	}
	return nil
}

func (client *HttpClient) Post(url string, contenType string, body string) error {
	resp, err := jammingHttp.Post(url, contenType, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("POST to %s: %v", url, err)
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
		return nil, "", fmt.Errorf("GET from %s: %v", url, geterr)
	}

	if !isStatusOk(resp) {
		return nil, "", fmt.Errorf("GET from %s: Status code is not OK (%d)", url, resp.StatusCode)
	}

	body, readerr := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if readerr != nil {
		return nil, "", fmt.Errorf("GET from %s: unable to read body: %v", url, readerr)
	}
	return resp.Header, string(body), nil
}

func isStatusOk(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode <= 299
}

func (client *HttpClient) Get(url string) (*http.Response, error) {
	return jammingHttp.Get(url)
}
