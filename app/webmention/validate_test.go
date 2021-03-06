package webmention

import (
	"errors"
	"net/http"
	"testing"

	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/mocks"
)

type httpReqMock struct {
	source string
	target string
}
type httpHeaderMock struct {
	contentType string
}

func (mock *httpHeaderMock) Get(key string) string {
	return mock.contentType
}
func (mock *httpReqMock) FormValue(key string) string {
	switch key {
	case "source":
		return mock.source
	case "target":
		return mock.target
	default:
		return ""
	}
}
func buildHttpReq(source string, target string) *httpReqMock {
	return &httpReqMock{
		source: source,
		target: target,
	}
}

var (
	config = &common.Config{
		BaseURL:                  "https://somecool.base.url/",
		Port:                     1337,
		Token:                    "miauwkes",
		AllowedWebmentionSources: []string{"brainbaking.com", "jefklakscodex.com"},
		Blacklist:                []string{"youtube.com"},
		Whitelist:                []string{"brainbaking.com"},
	}
)

func TestValidate(t *testing.T) {
	cases := []struct {
		label       string
		source      string
		target      string
		contentType string
		expected    bool
	}{
		{
			"is valid if encoded as form including the charset",
			"https://brainbaking.com/bla1",
			"https://jefklakscodex.com/bla",
			"application/x-www-form-urlencoded; charset=utf-8",
			true,
		},
		{
			"is valid if source and target https urls",
			"http://brainbaking.com/bla1",
			"http://jefklakscodex.com/bla",
			"application/x-www-form-urlencoded",
			true,
		},
		{
			"is NOT valid if target is a valid url but not form valid domain",
			"http://brainbaking.com/bla1",
			"http://brainthe.bake/jup",
			"application/x-www-form-urlencoded",
			false,
		},
		{
			"is NOT valid if source and target are the same urls",
			"http://brainbaking.com/bla1",
			"http://brainbaking.com/bla1",
			"application/x-www-form-urlencoded",
			false,
		},
		{
			"is NOT valid if source is not a valid url",
			"lolz",
			"http://brainbaking.com/bla1",
			"application/x-www-form-urlencoded",
			false,
		},
		{
			"is NOT valid if source is missing",
			"",
			"http://brainbaking.com/bla1",
			"application/x-www-form-urlencoded",
			false,
		},
		{
			"is NOT valid if target is missing",
			"http://brainbaking.com/bla1",
			"",
			"application/x-www-form-urlencoded",
			false,
		},
		{
			"is NOT valid if no valid encoded form",
			"http://brainbaking.com/bla1",
			"http://jefklakscodex.com/bla",
			"application/lolz",
			false,
		},
		{
			"is NOT valid if body is missing",
			"",
			"",
			"application/x-www-form-urlencoded",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			httpReq := buildHttpReq(tc.source, tc.target)
			httpHeader := &httpHeaderMock{contentType: tc.contentType}

			actual := validate(httpReq, httpHeader, config)
			if actual != tc.expected {
				t.Fatalf("got %v, want %v", actual, tc.expected)
			}
		})
	}
}

func TestIsValidTargetUrlFalseIfGetFails(t *testing.T) {
	client := &mocks.RestClientMock{
		GetFunc: func(url string) (*http.Response, error) {
			return nil, errors.New("whoops")
		},
	}
	result := isValidTargetUrl("failing", client)
	if result != false {
		t.Fatalf("expected to fail")
	}
}

func TestIsValidTargetUrlTrueIfGetSucceeds(t *testing.T) {
	client := &mocks.RestClientMock{
		GetFunc: func(url string) (*http.Response, error) {
			return nil, nil
		},
	}
	result := isValidTargetUrl("valid stuff!", client)
	if result != true {
		t.Fatalf("expected to succeed")
	}
}
