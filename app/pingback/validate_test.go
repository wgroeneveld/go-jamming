package pingback

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"brainbaking.com/go-jamming/common"
	"testing"
)

var conf *common.Config = &common.Config{
	AllowedWebmentionSources: []string{
		"brainbaking.com",
		"jefklakscodex.com",
	},
}

func TestValidate(t *testing.T) {
	cases := []struct {
		label string
		xml   string
		expected bool
	}{
		{
			"not valid if methodName is not pingback.ping",
			`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>ka.tsjing</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`,
			false,
		},
		{
			"not valid if less than two parameters",
			`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`,
			false,
		},
		{
			"not valid if more than two parameters",
			`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`,
			false,
		},
		{
			"not valid if target is not in trusted domains from config",
			`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>https://flashballz.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`,
			false,
		},
		{
			"not valid if target is not http(s)",
			`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>gemini://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`,
			false,
		},
		{
			"not valid if source is not http(s)",
			`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>gemini://cool.site</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`,
			false,
		},
		{
			"is valid if pingback.ping and two http(s) parameters of which target is trusted",
			`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`,
			true,
		},
		{
			"is not valid if source and target are the same urls",
			`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`,
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			var xmlObj XmlRPCMethodCall
			err := xml.Unmarshal([]byte(tc.xml), &xmlObj)
			assert.NoError(t, err, "XML invalid in test case")

			result := validate(&xmlObj, conf)
			assert.Equal(t, tc.expected, result)
		})
	}
}