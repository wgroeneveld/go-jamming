package send

import (
	"brainbaking.com/go-jamming/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

type CollectSuite struct {
	suite.Suite
	xml   string
	snder *Sender
}

func (s *CollectSuite) SetupTest() {
	file, _ := ioutil.ReadFile("../../../mocks/samplerss.xml")
	s.xml = string(file)
	s.snder = &Sender{
		Conf: &common.Config{
			Blacklist: []string{
				"youtube.com",
			},
		},
	}
}

func TestCollectSuite(t *testing.T) {
	suite.Run(t, new(CollectSuite))
}

func (s *CollectSuite) TestCollectUniqueHrefsFromHtmlShouldNotContainInlineLinks() {
	links := s.snder.collectUniqueHrefsFromHtml(`<html><body><a href="#inline">sup</a></body></html>`)
	assert.Empty(s.T(), links)
}

func (s *CollectSuite) TestCollectShouldNotContainHrefsFromBlockedDomains() {
	items, err := s.snder.Collect(s.xml, "https://brainbaking.com/notes/2021/03/09h15m17s30/")
	assert.NoError(s.T(), err)
	last := items[len(items)-1]
	assert.Equal(s.T(), "https://brainbaking.com/notes/2021/03/10h16m24s22/", last.link)
	assert.ElementsMatch(s.T(), []string{
		"https://dog.estate/@eli_oat",
		"https://twitter.com/olesovhcom/status/1369478732247932929",
		"/about",
	}, last.hrefs)
}

func (s *CollectSuite) TestCollectShouldNotContainHrefsThatPointToImages() {
	items, err := s.snder.Collect(s.xml, "https://brainbaking.com/notes/2021/03/13h12m44s29/")
	assert.NoError(s.T(), err)
	last := items[len(items)-1]
	// test case:
	// contains e.g. https://chat.brainbaking.com/media/6f8b72ca-9bfb-460b-9609-c4298a8cab2b/EuropeBattle%202021-03-14%2016-20-36-87.jpg
	assert.ElementsMatch(s.T(), []string{
		"/about",
	}, last.hrefs)
}

func (s *CollectSuite) TestCollectNothingIfNothingNewInFeed() {
	latestEntry := "https://brainbaking.com/notes/2021/03/16h17m07s14/"
	items, err := s.snder.Collect(s.xml, latestEntry)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(items))
}

func (s *CollectSuite) TestCollectLatestXLinksWhenARecentLinkParameterIsProvided() {
	items, err := s.snder.Collect(s.xml, "https://brainbaking.com/notes/2021/03/14h17m41s53/")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 3, len(items))

	last := items[len(items)-1]
	assert.Equal(s.T(), "https://brainbaking.com/notes/2021/03/15h14m43s49/", last.link)
	assert.ElementsMatch(s.T(), []string{
		"http://replit.com",
		"http://codepen.io",
		"https://kuleuven-diepenbeek.github.io/osc-course/ch1-c/intro/",
		"/about",
	}, last.hrefs)

}

func (s *CollectSuite) TestCollectEveryExternalLinkWithoutARecentLink() {
	items, err := s.snder.Collect(s.xml, "")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 141, len(items))

	first := items[0]
	assert.Equal(s.T(), "https://brainbaking.com/notes/2021/03/16h17m07s14/", first.link)
	assert.ElementsMatch(s.T(), []string{
		"https://fosstodon.org/@celia",
		"https://fosstodon.org/@kev",
		"/about",
	}, first.hrefs)

}
