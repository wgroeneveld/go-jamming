package common

import (
	"github.com/rs/zerolog/log"
	"time"
)

// https://labs.yulrizka.com/en/stubbing-time-dot-now-in-golang/
// None of the above are very appealing. For now, just use the lazy way.
var Now = time.Now

const (
	IsoFormat = "2006-01-02T15:04:05.000Z"
)

// TimeToIso converts time to ISO string format, up to seconds.
func TimeToIso(theTime time.Time) string {
	return theTime.Format(IsoFormat)
}

// IsoToTime converts an ISO time string into a time.Time object
// As produced by clients using day.js - e.g. 2021-04-09T15:51:43.732Z
func IsoToTime(since string) time.Time {
	if since == "" {
		return time.Time{}
	}
	t, err := time.Parse(IsoFormat, since)
	if err != nil {
		log.Warn().Str("time", since).Msg("Invalid ISO date, reverting to now()")
		return Now()
	}
	return t
}
