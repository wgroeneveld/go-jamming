package mf

import (
	"brainbaking.com/go-jamming/common"
	"strings"
	"time"
	"willnorris.com/go/microformats"
)

const (
	DateFormat = "2006-01-02T15:04:05"
)

type IndiewebAuthor struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type IndiewebData struct {
	Author       IndiewebAuthor `json:"author"`
	Name         string         `json:"name"`
	Content      string         `json:"content"`
	Published    string         `json:"published"`
	Url          string         `json:"url"`
	IndiewebType string         `json:"type"`
	Source       string         `json:"source"`
	Target       string         `json:"target"`
}

func PublishedNow(utcOffset int) string {
	return common.Now().UTC().Add(time.Duration(utcOffset) * time.Minute).Format("2006-01-02T15:04:05")
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
func Str(mf *microformats.Microformat, key string) string {
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

func Map(mf *microformats.Microformat, key string) map[string]string {
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

func Prop(mf *microformats.Microformat, key string) *microformats.Microformat {
	val := mf.Properties[key]
	if len(val) == 0 {
		return &microformats.Microformat{
			Properties: map[string][]interface{}{},
		}
	}
	return val[0].(*microformats.Microformat)
}

func DeterminePublishedDate(hEntry *microformats.Microformat, utcOffset int) string {
	publishedDate := Str(hEntry, "published")
	if publishedDate == "" {
		return PublishedNow(utcOffset)
	}
	return publishedDate
}

func DetermineAuthorName(hEntry *microformats.Microformat) string {
	authorName := Str(Prop(hEntry, "author"), "name")
	if authorName == "" {
		return Prop(hEntry, "author").Value
	}
	return authorName
}

func DetermineType(hEntry *microformats.Microformat) string {
	likeOf := Str(hEntry, "like-of")
	if likeOf != "" {
		return "like"
	}
	bookmarkOf := Str(hEntry, "bookmark-of")
	if bookmarkOf != "" {
		return "bookmark"
	}
	return "mention"
}

// Mastodon uids start with "tag:server", but we do want indieweb uids from other sources
func DetermineUrl(hEntry *microformats.Microformat, source string) string {
	uid := Str(hEntry, "uid")
	if uid != "" && strings.HasPrefix(uid, "http") {
		return uid
	}
	url := Str(hEntry, "url")
	if url != "" {
		return url
	}
	return source
}

func DetermineContent(hEntry *microformats.Microformat) string {
	bridgyTwitterContent := Str(hEntry, "bridgy-twitter-content")
	if bridgyTwitterContent != "" {
		return shorten(bridgyTwitterContent)
	}
	summary := Str(hEntry, "summary")
	if summary != "" {
		return shorten(summary)
	}
	contentEntry := Map(hEntry, "content")["value"]
	return shorten(contentEntry)
}
