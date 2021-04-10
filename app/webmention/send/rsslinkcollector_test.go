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
	xml string
}

func (s *CollectSuite) SetupTest() {
	file, _ := ioutil.ReadFile("../../../mocks/samplerss.xml")
	s.xml = string(file)
}

func TestCollectSuite(t *testing.T) {
	suite.Run(t, new(CollectSuite))
}

func (s *CollectSuite) TestCollectShouldNotContainHrefsFromBlockedDomains() {
	items, err := Collect(s.xml, common.IsoToTime("2021-03-10T00:00:00.000Z"))
	assert.NoError(s.T(), err)
	last := items[len(items)-1]
	assert.Equal(s.T(), "https://brainbaking.com/notes/2021/03/10h16m24s22/", last.link)
	/*
		assert.Equal(s.T(), []string{
			"https://dog.estate/@eli_oat",
			"https://twitter.com/olesovhcom/status/1369478732247932929",
			"/aobut",
		}, last.hrefs)

	*/
}

func (s *CollectSuite) TestCollectShouldNotContainHrefsThatPointToImages() {

}

func (s *CollectSuite) TestCollectIfTimeTagFoundInContextThatActsAsAnUpdateStamp() {

}

func (s *CollectSuite) TestCollectsNotIfTimeTagFoundInContextButStillOlderThanSince() {

}

func (s *CollectSuite) TestCollectNothingIfDateInFutureAndSinceNothingNewInFeed() {

}

func (s *CollectSuite) TestCollectLatestXLinksWhenASinceParameterIsProvided() {

}

func (s *CollectSuite) TestCollectEveryExternalLinkWithoutAValidSinceDate() {

}
