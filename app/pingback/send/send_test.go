package send

import (
	"github.com/stretchr/testify/assert"
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/mocks"
	"testing"
)

func TestSendPingbackToEndpoint(t *testing.T) {
	var capturedBody string
	sender := Sender{
		RestClient: &mocks.RestClientMock{
			PostFunc: func(url string, contentType string, body string) error {
				capturedBody = body
				return nil
			},
		},
	}
	expectedXml := `<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>src</string></value>
		</param>
		<param>
			<value><string>target</string></value>
		</param>
	</params>
</methodCall>`

	sender.SendPingbackToEndpoint("http://dingdong.com/pingback", mf.Mention{
		Source: "src",
		Target: "target",
	})
	assert.Equal(t, expectedXml, capturedBody)
}