package send

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/mocks"
	"github.com/stretchr/testify/assert"
	"net/url"
	"sync"
	"testing"
)

func TestSendMentionAsWebmention(t *testing.T) {
	passedFormValues := url.Values{}
	snder := Sender{
		RestClient: &mocks.RestClientMock{
			PostFormFunc: func(endpoint string, formValues url.Values) error {
				passedFormValues = formValues
				return nil
			},
		},
	}

	sendMentionAsWebmention(&snder, mf.Mention{
		Source: "mysource",
		Target: "mytarget",
	}, "someendpoint")

	assert.Equal(t, "mysource", passedFormValues.Get("source"))
	assert.Equal(t, "mytarget", passedFormValues.Get("target"))
}

func TestSendIntegrationTestCanSendBothWebmentionsAndPingbacks(t *testing.T) {
	posted := map[string]interface{}{}
	var lock = sync.RWMutex{}

	snder := Sender{
		Conf: common.Configure(),
		RestClient: &mocks.RestClientMock{
			GetBodyFunc: mocks.RelPathGetBodyFunc(t, "./../../../mocks/"),
			PostFunc: func(url string, contentType string, body string) error {
				lock.RLock()
				defer lock.RUnlock()
				posted[url] = body
				return nil
			},
			PostFormFunc: func(endpoint string, formValues url.Values) error {
				lock.RLock()
				defer lock.RUnlock()
				posted[endpoint] = formValues
				return nil
			},
		},
	}

	snder.Send("brainbaking.com", "2021-03-16T16:00:00.000Z")
	assert.Equal(t, 3, len(posted))

	wmPost1 := posted["http://aaronpk.example/webmention-endpoint-header"].(url.Values)
	assert.Equal(t, "https://brainbaking.com/notes/2021/03/16h17m07s14/", wmPost1.Get("source"))
	assert.Equal(t, "https://brainbaking.com/link-discover-test-multiple.html", wmPost1.Get("target"))

	wmPost2 := posted["http://aaronpk.example/pingback-endpoint-body"].(string)
	expectedPost2 := `<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://brainbaking.com/notes/2021/03/16h17m07s14/</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/pingback-discover-test-single.html</string></value>
		</param>
	</params>
</methodCall>`
	assert.Equal(t, expectedPost2, wmPost2)

	wmPost3 := posted["http://aaronpk.example/webmention-endpoint-body"].(url.Values)
	assert.Equal(t, "https://brainbaking.com/notes/2021/03/16h17m07s14/", wmPost3.Get("source"))
	assert.Equal(t, "https://brainbaking.com/link-discover-test-single.html", wmPost3.Get("target"))
}
