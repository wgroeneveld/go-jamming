package rest

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"io"
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

const (
	MaxBytes = 5000000 // 5 MiB
)

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

	ResponseAboveLimit = errors.New("response bigger than limit")
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

// GetBody issues a retryable GET request and returns the header, body string, and a possible error.
// It limits response sizes to MaxBytes and returns an error if status not between [200, 299].
func (client *HttpClient) GetBody(url string) (http.Header, string, error) {
	resp, geterr := client.Get(url)
	if geterr != nil {
		return nil, "", fmt.Errorf("GET from %s: %w", url, geterr)
	}

	if !isStatusOk(resp) {
		return nil, "", fmt.Errorf("GET from %s: Status code is not OK (%d)", url, resp.StatusCode)
	}

	body, readerr := readUntilMax(resp.Body, MaxBytes)
	defer resp.Body.Close()
	if readerr != nil {
		return nil, "", fmt.Errorf("GET from %s: unable to read body: %w", url, readerr)
	}
	return resp.Header, string(body), nil
}

func isStatusOk(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode <= 299
}

func (client *HttpClient) Get(url string) (*http.Response, error) {
	return jammingHttp.Get(url)
}

// readUntilMax is a duplicate of io.Read(). It behaves exactly the same.
// However, it will only read maxBytes bytes, exponentially chunked (as per append).
// Returns an error if it exceeds the limit.
func readUntilMax(r io.Reader, maxBytes int) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
		if len(b) > maxBytes {
			return nil, ResponseAboveLimit
		}
	}
}
