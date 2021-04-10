package send

import (
	"brainbaking.com/go-jamming/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
	"time"
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
			DisallowedWebmentionDomains: []string{
				"youtube.com",
			},
		},
	}
}

func TestCollectSuite(t *testing.T) {
	suite.Run(t, new(CollectSuite))
}

func (s *CollectSuite) TestCollectShouldNotContainHrefsFromBlockedDomains() {
	items, err := s.snder.Collect(s.xml, common.IsoToTime("2021-03-10T00:00:00.000Z"))
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
	items, err := s.snder.Collect(s.xml, common.IsoToTime("2021-03-14T00:00:00.000Z"))
	assert.NoError(s.T(), err)
	last := items[len(items)-1]
	// test case:
	// contains e.g. https://chat.brainbaking.com/media/6f8b72ca-9bfb-460b-9609-c4298a8cab2b/EuropeBattle%202021-03-14%2016-20-36-87.jpg
	assert.ElementsMatch(s.T(), []string{
		"/about",
	}, last.hrefs)
}

func (s *CollectSuite) TestCollectNothingIfDateInFutureAndSinceNothingNewInFeed() {
	items, err := s.snder.Collect(s.xml, time.Now().Add(time.Duration(600)*time.Hour))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(items))
}

func (s *CollectSuite) TestCollectLatestXLinksWhenASinceParameterIsProvided() {
	items, err := s.snder.Collect(s.xml, common.IsoToTime("2021-03-15T00:00:00.000Z"))
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

func (s *CollectSuite) TestCollectEveryExternalLinkWithoutAValidSinceDate() {
	// no valid since date = zero time passed.
	items, err := s.snder.Collect(s.xml, time.Time{})
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
