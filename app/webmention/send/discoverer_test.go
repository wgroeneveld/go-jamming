package send

import (
	"brainbaking.com/go-jamming/mocks"
	"brainbaking.com/go-jamming/rest"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestDiscoverMentionEndpointE2EWithRedirect(t *testing.T) {
	t.Skip("Skipping TestDiscoverE2EWithRedirect, webmention.rocks is slow.")
	var sender = &Sender{
		RestClient: &rest.HttpClient{},
	}

	link, wmType := sender.discoverMentionEndpoint("https://webmention.rocks/test/23/page")
	assert.Equal(t, typeWebmention, wmType)
	expectedUrl := "https://webmention.rocks/test/23/page/webmention-endpoint/"
	assert.Truef(t, strings.HasPrefix(link, expectedUrl), "should start with %s, but was %s", expectedUrl, link)
}

func TestDisccoverRssFeedPrefersFirstEntriesOverLater(t *testing.T) {
	var snder = &Sender{
		RestClient: &mocks.RestClientMock{
			HeadFunc: func(s string) (*http.Response, error) {
				if strings.HasSuffix(s, "/index.xml") {
					return &http.Response{
						Header: map[string][]string{
							"Content-Type": {"text/xml"},
						},
						StatusCode: 200,
					}, nil
				}
				return nil, fmt.Errorf("BOOM")
			},
		},
	}

	feed, err := snder.discoverRssFeed("blah.com")
	assert.NoError(t, err)
	assert.Equal(t, "https://blah.com/all/index.xml", feed)
}

func TestDiscoverRssFeedNoneFoundReturnsError(t *testing.T) {
	var snder = &Sender{
		RestClient: &mocks.RestClientMock{
			HeadFunc: func(s string) (*http.Response, error) {
				return nil, fmt.Errorf("BOOM")
			},
		},
	}

	_, err := snder.discoverRssFeed("blah.com")
	assert.Error(t, err)
}
func TestDiscoverRssFeedFirstNotXmlReturnsSecondWorkingOne(t *testing.T) {
	var snder = &Sender{
		RestClient: &mocks.RestClientMock{
			HeadFunc: func(s string) (*http.Response, error) {
				if strings.HasSuffix(s, "/all/index.xml") {
					return &http.Response{
						Header: map[string][]string{
							"Content-Type": {"text/html"},
						},
						StatusCode: 200,
					}, nil
				}
				if strings.HasSuffix(s, "/feed") {
					return &http.Response{
						Header: map[string][]string{
							"Content-Type": {"text/xml"},
						},
						StatusCode: 200,
					}, nil
				}
				return nil, fmt.Errorf("BOOM")
			},
		},
	}

	feed, err := snder.discoverRssFeed("blah.com")
	assert.NoError(t, err)
	assert.Equal(t, "https://blah.com/feed", feed)
}

func TestDiscoverMentionEndpoint(t *testing.T) {
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
			"https://brainbaking.com/link-discover-test-none.html",
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
			"webmentions: https://webmention.rocks/test/1 relative path in header",
			"https://brainbaking.com/webmention-rocks-1.html",
			"https://brainbaking.com/test/1/webmention?head=true",
			typeWebmention,
		},
		{
			"webmentions: https://webmention.rocks/test/11 prefer links in the header even if in head or body also present",
			"https://brainbaking.com/webmention-rocks-11.html",
			"https://brainbaking.com/test/11/webmention",
			typeWebmention,
		},
		{
			"webmentions: https://webmention.rocks/test/15 empty link rel means it is its own endpoint",
			"https://brainbaking.com/webmention-rocks-15.html",
			"https://brainbaking.com/webmention-rocks-15.html",
			typeWebmention,
		},
		{
			"webmentions: https://webmention.rocks/test/18 discover link if multiple present in multiple headers",
			"https://brainbaking.com/webmention-rocks-18.html",
			"https://webmention.rocks/test/18/webmention?head=true",
			typeWebmention,
		},
		{
			"webmentions: https://webmention.rocks/test/19 discover link if multiple present comma-separated in single header",
			"https://brainbaking.com/webmention-rocks-19.html",
			"https://webmention.rocks/test/19/webmention?head=true",
			typeWebmention,
		},
		{
			"webmentions: https://webmention.rocks/test/20 discover link in body href if header is empty",
			"https://brainbaking.com/webmention-rocks-20.html",
			"https://brainbaking.com/test/20/webmention",
			typeWebmention,
		},
		{
			"webmentions: https://webmention.rocks/test/22 discover link relative to the page instead of the domain",
			"https://brainbaking.com/blank/webmention-rocks-22.html",
			"https://brainbaking.com/blank/22/webmention",
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
			link, mentionType := sender.discoverMentionEndpoint(tc.url)
			assert.Equal(t, tc.expectedLink, link)
			assert.Equal(t, tc.expectedType, mentionType)
		})
	}
}
