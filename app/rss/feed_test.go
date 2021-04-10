package rss

import (
	"brainbaking.com/go-jamming/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPubDateAsTimeIncorrectRevertsToNow(t *testing.T) {
	common.Now = func() time.Time {
		return time.Date(2020, time.January, 1, 12, 30, 0, 0, time.UTC)
	}
	itm := Item{
		PubDate: "frutselbolletjes",
	}

	theTime := itm.PubDateAsTime()
	assert.Equal(t, 2020, theTime.Year())
	assert.Equal(t, time.January, theTime.Month())
}

func TestPubDateAsTime(t *testing.T) {
	itm := Item{
		PubDate: "Tue, 16 Mar 2021 17:07:14 +0000",
	}
	theTime := itm.PubDateAsTime()
	assert.Equal(t, 2021, theTime.Year())
	assert.Equal(t, time.March, theTime.Month())
	assert.Equal(t, 16, theTime.Day())
	assert.Equal(t, 17, theTime.Hour())
	assert.Equal(t, 7, theTime.Minute())
	assert.Equal(t, 14, theTime.Second())
}
