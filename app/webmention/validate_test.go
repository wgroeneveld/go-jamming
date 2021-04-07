
package webmention

import (
	"testing"

	"github.com/wgroeneveld/go-jamming/common"
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
	case "source": return mock.source
	case "target": return mock.target
	default: return ""
	}
}
func buildHttpReq(source string, target string) *httpReqMock {
	return &httpReqMock{
		source: source,
		target: target,
	}
}

var config = common.Configure()

func TestValidate(t *testing.T) {
	cases := []struct {
		label string
		source string
		target string
		contentType string
		expected bool
	} {
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
			httpHeader := &httpHeaderMock{ contentType: tc.contentType }

			actual := validate(httpReq, httpHeader, config)
			if actual != tc.expected {
				t.Fatalf("got %v, want %v", actual, tc.expected)
			}
		})
	}

	
}
