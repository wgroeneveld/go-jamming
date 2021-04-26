package mf

import (
	"brainbaking.com/go-jamming/common"
	"fmt"
	"strings"
	"time"
	"willnorris.com/go/microformats"
)

const (
	dateFormatWithTimeZone               = "2006-01-02T15:04:05-07:00"
	dateFormatWithAbsoluteTimeZone       = "2006-01-02T15:04:05-0700"
	dateFormatWithTimeZoneSuffixed       = "2006-01-02T15:04:05.000Z"
	dateFormatWithoutTimeZone            = "2006-01-02T15:04:05"
	dateFormatWithSecondsWithoutTimeZone = "2006-01-02T15:04:05.00Z"
	dateFormatWithoutTime                = "2006-01-02"
	Anonymous                            = "anonymous"
)

var (
	supportedFormats = []string{
		dateFormatWithTimeZone,
		dateFormatWithAbsoluteTimeZone,
		dateFormatWithTimeZoneSuffixed,
		dateFormatWithSecondsWithoutTimeZone,
		dateFormatWithoutTimeZone,
		dateFormatWithoutTime,
	}
)

type IndiewebAuthor struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (ia *IndiewebAuthor) Anonymize() {
	ia.Picture = fmt.Sprintf("/pictures/%s", Anonymous)
}

type IndiewebDataResult struct {
	Status string          `json:"status"`
	Data   []*IndiewebData `json:"json"`
}

func ResultFailure(data []*IndiewebData) IndiewebDataResult {
	return IndiewebDataResult{
		Status: "failure",
		Data:   data,
	}
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

func (id *IndiewebData) AsMention() Mention {
	return Mention{
		Source: id.Source,
		Target: id.Target,
	}
}

func (id *IndiewebData) IsEmpty() bool {
	return id.Url == ""
}

func PublishedNow(zone *time.Location) string {
	return common.Now().UTC().In(zone).Format(dateFormatWithTimeZone)
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
			return ""
		}
		return valMap["value"]
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
	return hItemType(data, "h-entry")
}

func HCard(data *microformats.Data) *microformats.Microformat {
	return hItemType(data, "h-card")
}

func hItemType(data *microformats.Data, hType string) *microformats.Microformat {
	for _, itm := range data.Items {
		if common.Includes(itm.Type, hType) {
			return itm
		}
	}
	return nil
}

func mfEmpty() *microformats.Microformat {
	return &microformats.Microformat{
		Properties: map[string][]interface{}{},
	}
}

func Prop(mf *microformats.Microformat, key string) *microformats.Microformat {
	val := mf.Properties[key]
	if len(val) == 0 {
		return mfEmpty()
	}
	for i := range val {
		conv, ok := val[i].(*microformats.Microformat)
		if ok {
			return conv
		}
	}
	return mfEmpty()
}

func Published(hEntry *microformats.Microformat, zone *time.Location) string {
	publishedDate := Str(hEntry, "published")
	if publishedDate == "" {
		return PublishedNow(zone)
	}

	for _, format := range supportedFormats {
		formatted, err := time.Parse(format, publishedDate)
		if err != nil {
			continue
		}
		return formatted.Format(dateFormatWithTimeZone)
	}

	return PublishedNow(zone)
}

func NewAuthor(hEntry *microformats.Microformat, hCard *microformats.Microformat) IndiewebAuthor {
	name := DetermineAuthorName(hEntry)
	if name == "" {
		name = DetermineAuthorName(hCard)
	}
	picture := DetermineAuthorPhoto(hEntry)
	if picture == "" {
		picture = DetermineAuthorPhoto(hCard)
	}
	return IndiewebAuthor{
		Picture: picture,
		Name:    name,
	}
}

func DetermineAuthorPhoto(hEntry *microformats.Microformat) string {
	photo := Str(Prop(hEntry, "author"), "photo")
	if photo == "" {
		photo = Str(hEntry, "photo")
	}
	return photo
}

func DetermineAuthorName(hEntry *microformats.Microformat) string {
	authorName := Str(Prop(hEntry, "author"), "name")
	if authorName == "" {
		authorName = Prop(hEntry, "author").Value
	}
	if authorName == "" {
		authorName = Str(hEntry, "author")
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
