package pingback

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

// See https://www.hixie.ch/specs/pingback/pingback#refsXMLRPC
func TestMarshallValidXMLRPC(t *testing.T) {
	xmlString := `<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://brainbaking.com/kristien.html</string></value>
		</param>
		<param>
			<value><string>https://kristienthoelen.be/2021/03/22/de-stadia-van-een-burn-out-in-welk-stadium-zit-jij/</string></value>
		</param>
	</params>
</methodCall>`
	var rpc XmlRPCMethodCall
	err := xml.Unmarshal([]byte(xmlString), &rpc)

	assert.NoError(t, err)
	assert.Equal(t, "pingback.ping", rpc.MethodName)
	assert.Equal(t, "https://brainbaking.com/kristien.html", rpc.Params.Parameters[0].Value.String)
	assert.Equal(t, "https://kristienthoelen.be/2021/03/22/de-stadia-van-een-burn-out-in-welk-stadium-zit-jij/", rpc.Params.Parameters[1].Value.String)
}
