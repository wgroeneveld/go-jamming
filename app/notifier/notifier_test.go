package notifier

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildNotification(t *testing.T) {
	wm := mf.Mention{
		Source: "https://brainbaking.com/valid-indieweb-source.html",
		Target: "https://brainbaking.com/valid-indieweb-target.html",
	}
	cnf := &common.Config{
		AllowedWebmentionSources: []string{
			"brainbaking.com",
		},
		BaseURL:   "https://jam.brainbaking.com/",
		Token:     "mytoken",
		Blacklist: []string{},
		Whitelist: []string{},
	}

	expected := `Hi admin, 

,A webmention was received: 
Source https://brainbaking.com/valid-indieweb-source.html, Target https://brainbaking.com/valid-indieweb-target.html
Content: somecontent

Accept? https://jam.brainbaking.com/admin/approve/mytoken/19d462ddff3c3322c662dac3461324bb:brainbaking.com
Reject? https://jam.brainbaking.com/admin/reject/mytoken/19d462ddff3c3322c662dac3461324bb:brainbaking.com
Cheerio, your go-jammin' thing.`

	result := BuildNotification(wm, &mf.IndiewebData{Content: "somecontent"}, cnf)
	assert.Equal(t, result, expected)
}
