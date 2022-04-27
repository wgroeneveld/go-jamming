package mf

import (
	"brainbaking.com/go-jamming/common"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"willnorris.com/go/microformats"
)

func TestResultSuccessNonEmpty(t *testing.T) {
	arr := make([]*IndiewebData, 1)
	arr[0] = &IndiewebData{Author: IndiewebAuthor{Name: "Jaak"}}

	data := ResultSuccess(arr)
	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)

	expected := `{"status":"success","json":[{"author":{"name":"Jaak","picture":""},"name":"","content":"","published":"","url":"","type":"","source":"","target":""}]}`
	assert.Equal(t, expected, string(jsonData))
}

func TestResultSuccessEmptyEncodesAsEmptyJSONArray(t *testing.T) {
	data := ResultSuccess(nil)
	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)

	expected := `{"status":"success","json":[]}`
	assert.Equal(t, expected, string(jsonData))
}

func TestResultFailureNonEmpty(t *testing.T) {
	arr := make([]*IndiewebData, 1)
	arr[0] = &IndiewebData{Author: IndiewebAuthor{Name: "Jaak"}}

	data := ResultFailure(arr)
	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)

	expected := `{"status":"failure","json":[{"author":{"name":"Jaak","picture":""},"name":"","content":"","published":"","url":"","type":"","source":"","target":""}]}`
	assert.Equal(t, expected, string(jsonData))
}

func TestResultFailureEmptyEncodesAsEmptyJSONArray(t *testing.T) {
	data := ResultFailure(nil)
	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)

	expected := `{"status":"failure","json":[]}`
	assert.Equal(t, expected, string(jsonData))
}

func TestPublished(t *testing.T) {
	common.Now = func() time.Time {
		return time.Date(2020, time.January, 1, 12, 30, 0, 0, time.UTC)
	}
	nowString := "2020-01-01T12:30:00+00:00"
	defer func() {
		common.Now = time.Now
	}()
	cases := []struct {
		label        string
		raw          string
		expectedTime string
	}{
		{
			"Converts published date in RFC3339 ISO8601 indieweb datetime format with timezone",
			"2021-04-25T11:24:48+02:00",
			"2021-04-25T11:24:48+02:00",
		},
		{
			"Converts published date in RFC3339 ISO8601 indieweb datetime format with absolute timezone",
			"2021-04-25T11:24:48+0200",
			"2021-04-25T11:24:48+02:00",
		},
		{
			"Converts published date in RFC3339 ISO8601 indieweb datetime format with timezone suffixed with Z",
			"2021-03-02T16:17:18.000Z",
			"2021-03-02T16:17:18+00:00",
		},
		{
			"Converts published date in RFC3339 ISO8601 indieweb datetime format without timezone",
			"2021-04-25T11:24:48",
			"2021-04-25T11:24:48+00:00",
		},
		{
			"Converts published date in RFC3339 ISO8601 indieweb datetime format without time",
			"2021-04-25",
			"2021-04-25T00:00:00+00:00",
		},
		{
			"Returns current UTC date if property with correct timezone not found",
			"",
			nowString,
		},
		{
			"Reverts to current UTC date if not in correct ISO8601 datetime format",
			"26 April 2021",
			nowString,
		},
		{
			"https://www.ietf.org/rfc/rfc3339.txt example 1",
			"1985-04-12T23:20:50.52Z",
			"1985-04-12T23:20:50+00:00",
		},
		{
			"https://www.ietf.org/rfc/rfc3339.txt example 2",
			"1996-12-19T16:39:57-08:00",
			"1996-12-19T16:39:57-08:00",
		},
		{
			"https://www.ietf.org/rfc/rfc3339.txt example 3 explicitly not implemented",
			"1990-12-31T23:59:60Z",
			nowString,
		},
		{
			"https://www.ietf.org/rfc/rfc3339.txt example 4 explicitly not implemented",
			"1990-12-31T15:59:60-08:00",
			nowString,
		},
		{
			"https://www.ietf.org/rfc/rfc3339.txt example 5 with seconds ignored",
			"1937-01-01T12:00:27.87+00:20",
			"1937-01-01T12:00:27+00:20",
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			props := map[string][]interface{}{}
			props["published"] = []interface{}{
				tc.raw,
			}
			theTime := Published(&microformats.Microformat{
				Properties: props,
			})

			assert.Equal(t, tc.expectedTime, theTime)
		})

	}
}
