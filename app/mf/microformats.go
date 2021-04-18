package mf

import (
	"brainbaking.com/go-jamming/common"
	"encoding/json"
	"io/ioutil"
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

type IndiewebDataResult struct {
	Status string          `json:"status"`
	Data   []*IndiewebData `json:"json"`
}

func ResultSuccess(data []*IndiewebData) IndiewebDataResult {
	return IndiewebDataResult{
		Status: "success",
		Data:   data,
	}
}

type IndiewebData struct {
	Author       IndiewebAuthor `json:"author"`
	Name         string         `json:"name"`
	Content      string         `json:"content"`
	Published    string         `json:"published"`
	Url          string         `json:"url"`
	IndiewebType MfType         `json:"type"`
	Source       string         `json:"source"`
	Target       string         `json:"target"`
}

func (id *IndiewebData) IsEmpty() bool {
	return id.Url == ""
}

// RequireFromFile converts the file JSON contents into the indieweb struct.
// This ignores read and marshall errors and returns an emtpy struct instead.
func RequireFromFile(file string) *IndiewebData {
	indiewebData := &IndiewebData{}
	data, _ := ioutil.ReadFile(file)
	json.Unmarshal(data, indiewebData)
	return indiewebData
}

func PublishedNow(utcOffset int) string {
	return common.Now().UTC().Add(time.Duration(utcOffset) * time.Minute).Format("2006-01-02T15:04:05")
}

func shorten(txt string) string {
	if len(txt) <= 250 {
		return txt
	}
	return txt[:250] + "..."
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

func HEntry(data *microformats.Data) *microformats.Microformat {
	for _, itm := range data.Items {
		if common.Includes(itm.Type, "h-entry") {
			return itm
		}
	}
	return nil
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

func Published(hEntry *microformats.Microformat, utcOffset int) string {
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

type MfType string

const (
	TypeLike     MfType = "like"
	TypeBookmark MfType = "bookmark"
	TypeMention  MfType = "mention"
)

func Type(hEntry *microformats.Microformat) MfType {
	likeOf := Str(hEntry, "like-of")
	if likeOf != "" {
		return TypeLike
	}
	bookmarkOf := Str(hEntry, "bookmark-of")
	if bookmarkOf != "" {
		return TypeBookmark
	}
	return TypeMention
}

// Mastodon uids start with "tag:server", but we do want indieweb uids from other sources
func Url(hEntry *microformats.Microformat, source string) string {
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

func Content(hEntry *microformats.Microformat) string {
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
