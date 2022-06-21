package common

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TimeSuite struct {
	suite.Suite
	nowtime time.Time
}

func (s *TimeSuite) SetupTest() {
	s.nowtime = time.Date(2020, time.January, 1, 12, 30, 0, 0, time.UTC)
	Now = func() time.Time {
		return s.nowtime
	}
}

func TestSendSuite(t *testing.T) {
	suite.Run(t, new(TimeSuite))
}

func (s *TimeSuite) TestTimeToIso() {
	theTime := time.Date(2021, time.March, 9, 15, 51, 43, 732, time.UTC)
	expected := "2021-03-09T15:51:43.000Z"
	actual := TimeToIso(theTime)

	assert.Equal(s.T(), expected, actual)
}

func (s *TimeSuite) TestIsoToTimeInISOString() {
	expectedtime := time.Date(2021, time.March, 9, 15, 51, 43, 732, time.UTC)
	since := IsoToTime("2021-03-09T15:51:43.732Z")
	assert.Equal(s.T(), expectedtime.Year(), since.Year())
	assert.Equal(s.T(), expectedtime.Month(), since.Month())
	assert.Equal(s.T(), expectedtime.Day(), since.Day())
	assert.Equal(s.T(), expectedtime.Hour(), since.Hour())
	assert.Equal(s.T(), expectedtime.Minute(), since.Minute())
	assert.Equal(s.T(), expectedtime.Second(), since.Second())
}

func (s *TimeSuite) TestIsoToTimeInvalidStringReturnsZeroTime() {
	since := IsoToTime("woef ik ben een hondje")
	assert.True(s.T(), since.IsZero())
}

func (s *TimeSuite) TestIsoToTimeEmptyReturnsZeroTime() {
	since := IsoToTime("")
	assert.True(s.T(), since.IsZero())
}
