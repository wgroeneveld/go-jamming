package webmention

import (
	"strings"
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

func determinePublishedDate(hEntry *microformats.Microformat, utcOffset int) string {
	publishedDate := mfStr(hEntry, "published")
	if publishedDate == "" {
		return publishedNow(utcOffset)
	}
	return publishedDate
}

func determineAuthorName(hEntry *microformats.Microformat) string {
	authorName := mfStr(mfProp(hEntry, "author"), "name")
	if authorName == "" {
		return mfProp(hEntry, "author").Value
	}
	return authorName
}

func determineMfType(hEntry *microformats.Microformat) string {
	likeOf := mfStr(hEntry, "like-of")
	if likeOf != "" {
		return "like"
	}
	bookmarkOf := mfStr(hEntry, "bookmark-of")
	if bookmarkOf != "" {
		return "bookmark"
	}
	return "mention"
}

// Mastodon uids start with "tag:server", but we do want indieweb uids from other sources
func determineUrl(hEntry *microformats.Microformat, source string) string {
	uid := mfStr(hEntry, "uid")
	if uid != "" && strings.HasPrefix(uid, "http") {
		return uid
	}
	url := mfStr(hEntry, "url")
	if url != "" {
		return url
	}
	return source
}

func determineContent(hEntry *microformats.Microformat) string {
	bridgyTwitterContent := mfStr(hEntry, "bridgy-twitter-content")
	if bridgyTwitterContent != "" {
		return shorten(bridgyTwitterContent)
	}
	summary := mfStr(hEntry, "summary")
	if summary != "" {
		return shorten(summary)
	}
	contentEntry := mfMap(hEntry, "content")["value"]
	return shorten(contentEntry)
}