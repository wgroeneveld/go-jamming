package send

import (
	"brainbaking.com/go-jamming/common"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestCollectUniqueHrefsFromHtml(t *testing.T) {
	cases := []struct {
		label         string
		html          string
		expectedLinks []string
	}{
		{
			"should not contain inline links",
			`<html><body><a href="#inline">sup</a></body></html>`,
			[]string{},
		},
		{
			"should not collect blacklisted links",
			`<html><body><a href="https://www.blacklisted.com/wowo.html">sup</a> and also <a href="/dinges">dinges</a>!</body></html>`,
			[]string{
				"/dinges",
			},
		},
		{
			"should not collect hrefs from <link/> tags, only from <a/> ones",
			`<html><head><link rel="stylesheet" href="/style.css"></head><body><a href="/dinges">dinges</a>!</body></html>`,
			[]string{
				"/dinges",
			},
		},
		{
			"should collect even if href is not the first attribute of an <a> tag",
			`<html><body><a style="cool" target="_blank" href="/one">one</a> and <a target="_blank" href="/two">two</a> and <a href="/three">three</a></body></html>`,
			[]string{
				"/one",
				"/two",
				"/three",
			},
		},
		{
			"should collect case insensitive",
			`<html><body><A href="/one">one</A> and <a href="/two">two</a> and <a HREF="/three">three</a></body></html>`,
			[]string{
				"/one",
				"/two",
				"/three",
			},
		},
		{
			"should not collect zips or ZIPs or gifs or GIFS",
			`<a href="/cool.gif">cool gif</a> and <a href="/more-cool.GIF">more-cool gif</a> and here's a zip: <a href="baf.ZIP">baf</a> or <a href="boef.zip">boef.zip</a>??'`,
			[]string{},
		},
	}

	s := &Sender{
		Conf: &common.Config{
			Blacklist: []string{
				"blacklisted.com",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			result := s.collectUniqueHrefsFromHtml(tc.html)
			assert.ElementsMatch(t, tc.expectedLinks, result)
		})
	}
}

func TestCollect(t *testing.T) {
	file, _ := ioutil.ReadFile("../../../mocks/samplerss.xml")
	snder := &Sender{
		Conf: &common.Config{
			Blacklist: []string{
				"youtube.com",
			},
		},
	}

	cases := []struct {
		label             string
		lastsentlink      string
		expectedRssItems  int
		expectedLastLinks []string
	}{
		{
			"should not contain hrefs from blocked domains",
			"https://brainbaking.com/notes/2021/03/09h15m17s30/",
			10,
			[]string{
				"https://dog.estate/@eli_oat",
				"https://twitter.com/olesovhcom/status/1369478732247932929",
				"/about",
			},
		},
		{
			// test case: contains e.g. https://chat.brainbaking.com/media/6f8b72ca-9bfb-460b-9609-c4298a8cab2b/EuropeBattle%202021-03-14%2016-20-36-87.jpg
			"should not contain hrefs that point to images",
			"https://brainbaking.com/notes/2021/03/13h12m44s29/",
			4,
			[]string{
				"/about",
			},
		},
		{
			"collects nothing if nothing new in feed",
			"https://brainbaking.com/notes/2021/03/16h17m07s14/",
			0,
			[]string{},
		},
		{
			"collect latest X links when a recent link parameter is provided",
			"https://brainbaking.com/notes/2021/03/14h17m41s53/",
			3,
			[]string{
				"http://replit.com",
				"http://codepen.io",
				"https://kuleuven-diepenbeek.github.io/osc-course/ch1-c/intro/",
				"/about",
			},
		},
		{
			"collect every external link without a recent link",
			"",
			141,
			[]string{
				"/notes/index.xml",
				"/archives",
				"/categories/hardware/index.xml",
				"/about",
				"https://netnewswire.com/",
				"/index.xml",
				"brainbaking.com",
				"/post/index.xml",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			items, err := snder.Collect(string(file), tc.lastsentlink)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRssItems, len(items))

			if tc.expectedRssItems > 0 {
				last := items[len(items)-1]
				assert.ElementsMatch(t, tc.expectedLastLinks, last.hrefs)
			}
		})
	}
}
