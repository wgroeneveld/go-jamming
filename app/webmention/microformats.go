package webmention

import (
	"time"
	"willnorris.com/go/microformats"
)

const (
	DateFormat = "2006-01-02T15:04:05"
)

type indiewebAuthor struct {
	Name string				`json:"name"`
	Picture string			`json:"picture"`
}

type indiewebData struct {
	Author indiewebAuthor	`json:"author"`
	Name string				`json:"name"`
	Content string			`json:"content"`
	Published string		`json:"published"`
	Url string				`json:"url"`
	IndiewebType string 	`json:"type"`
	Source string			`json:"source"`
	Target string			`json:"target"`
}

var now = time.Now
func publishedNow(utcOffset int) string {
	return now().UTC().Add(time.Duration(utcOffset) * time.Minute).Format("2006-01-02T15:04:05")
}

func shorten(txt string) string {
	if len(txt) <= 250 {
		return txt
	}
	return txt[0:250] + "..."
}

// Go stuff: entry.Properties["name"][0].(string),
// JS stuff: hEntry.properties?.name?.[0]
// The problem: convoluted syntax and no optional chaining!
func mfStr(mf *microformats.Microformat, key string) string {
	val := mf.Properties[key]
	if len(val) == 0 {
		return ""
	}

	str, ok := val[0].(string)
	if !ok {
		// in very weird cases, it could be a map holding a value, like in mf2's "photo"
		valMap, ok2 := val[0].(map[string]string)
		if !ok2 {
			str = ""
		}
		str = valMap["value"]
	}

	return str
}

func mfMap(mf *microformats.Microformat, key string) map[string]string {
	val := mf.Properties[key]
	if len(val) == 0 {
		return map[string]string{}
	}
	mapVal, ok := val[0].(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return mapVal
}

func mfProp(mf *microformats.Microformat, key string) *microformats.Microformat {
	val := mf.Properties[key]
	if len(val) == 0 {
		return &microformats.Microformat{
			Properties: map[string][]interface{}{},
		}
	}
	return val[0].(*microformats.Microformat)
}
