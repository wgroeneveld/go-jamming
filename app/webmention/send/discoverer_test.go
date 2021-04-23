package send

import (
	"brainbaking.com/go-jamming/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiscover(t *testing.T) {
	var sender = &Sender{
		RestClient: &mocks.RestClientMock{
			GetBodyFunc: mocks.RelPathGetBodyFunc("../../../mocks/"),
		},
	}

	cases := []struct {
		label        string
		url          string
		expectedLink string
		expectedType string
	}{
		{
			"discover 'unknown' if no link is present",
			"https://brainbaking.com/link-discover-test-none.html",
			"",
			typeUnknown,
		},
		{
			"prefer webmentions over pingbacks if both links are present",
			"https://brainbaking.com/link-discover-bothtypes.html",
			"http://aaronpk.example/webmention-endpoint",
			typeWebmention,
		},
		{
			"pingbacks: discover link if present in header",
			"https://brainbaking.com/pingback-discover-test.html",
			"http://aaronpk.example/pingback-endpoint",
			typePingback,
		},
		{
			"pingbacks: discover link if sole entry somewhere in html",
			"https://brainbaking.com/pingback-discover-test-single.html",
			"http://aaronpk.example/pingback-endpoint-body",
			typePingback,
		},
		{
			"pingbacks: use link in header if multiple present in html",
			"https://brainbaking.com/pingback-discover-test-multiple.html",
			"http://aaronpk.example/pingback-endpoint-header",
			typePingback,
		},
		{
			"webmentions: discover link if present in header",
			"https://brainbaking.com/link-discover-test.html",
			"http://aaronpk.example/webmention-endpoint",
			typeWebmention,
		},
		{
			"webmentions: discover link if sole entry somewhere in html",
			"https://brainbaking.com/link-discover-test-single.html",
			"http://aaronpk.example/webmention-endpoint-body",
			typeWebmention,
		},
		{
			"webmentions: use link in header if multiple present in html",
			"https://brainbaking.com/link-discover-test-multiple.html",
			"http://aaronpk.example/webmention-endpoint-header",
			typeWebmention,
		},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			link, mentionType := sender.discover(tc.url)
			assert.Equal(t, tc.expectedLink, link)
			assert.Equal(t, tc.expectedType, mentionType)
		})
	}
}
