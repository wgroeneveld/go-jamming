package common

import (
	"github.com/rs/zerolog/log"
	"time"
)

// https://labs.yulrizka.com/en/stubbing-time-dot-now-in-golang/
// None of the above are very appealing. For now, just use the lazy way.
var Now = time.Now

// since should be in ISO String format, as produced by clients using day.js - e.g. 2021-04-09T15:51:43.732Z
func IsoToTime(since string) time.Time {
	if since == "" {
		return time.Time{}
	}
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, since)
	if err != nil {
		log.Warn().Str("time", since).Msg("Invalid ISO date, reverting to now()")
		return Now()
	}
	return t
}
