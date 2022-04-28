package notifier

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildReceivedMsgDoesNotContainApproveLink(t *testing.T) {
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

	result := buildReceivedMsg(wm, &mf.IndiewebData{Content: "somecontent"}, cnf)
	assert.Contains(t, result, `Webmention approved from`)
	assert.Contains(t, result, `<em>Source:</em> <a href="https://brainbaking.com/valid-indieweb-source.html">https://brainbaking.com/valid-indieweb-source.html</a><br/>`)
	assert.Contains(t, result, `<em>Target:</em> <a href="https://brainbaking.com/valid-indieweb-target.html">https://brainbaking.com/valid-indieweb-target.html</a><br/>`)
	assert.NotContains(t, result, `<a href="https://jam.brainbaking.com/admin/approve/mytoken/19d462ddff3c3322c662dac3461324bb:brainbaking.com`)
}

func TestBuildInModerationMsgContainsApproveLink(t *testing.T) {
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

	result := buildInModerationMsg(wm, &mf.IndiewebData{Content: "somecontent"}, cnf)
	assert.Contains(t, result, `Webmention in moderation from`)
	assert.Contains(t, result, `<em>Source:</em> <a href="https://brainbaking.com/valid-indieweb-source.html">https://brainbaking.com/valid-indieweb-source.html</a><br/>`)
	assert.Contains(t, result, `<em>Target:</em> <a href="https://brainbaking.com/valid-indieweb-target.html">https://brainbaking.com/valid-indieweb-target.html</a><br/>`)
	assert.Contains(t, result, `<a href="https://jam.brainbaking.com/admin/approve/mytoken/19d462ddff3c3322c662dac3461324bb:brainbaking.com`)
}
