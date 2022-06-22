package external

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
	"encoding/json"
)

/*
An example webmention.io JSON mention object:
{
	"source": "https://brid.gy/like/twitter/iamchrisburnell/1298550501307486208/252048752",
	"verified": true,
	"verified_date": "2022-06-10T08:20:16+00:00",
	"id": 1412862,
	"private": false,
	"data": {
		"author": {
			"name": "Felipe Sere",
			"url": "https://twitter.com/felipesere",
			"photo": "https://webmention.io/avatar/pbs.twimg.com/ca88219f9ffb14d73bfb2bf88450cadb355bdaf56773947f39c752af3883fecf.jpg"
		},
		"url": "https://twitter.com/iamchrisburnell/status/1298550501307486208#favorited-by-252048752",
		"name": null,
		"content": null,
		"published": null,
		"published_ts": null
	},
	"activity": {
		"type": "like",
		"sentence": "Felipe Sere favorited a tweet https://chrisburnell.com/bowhead/",
		"sentence_html": "<a href=\"https://twitter.com/felipesere\">Felipe Sere</a> favorited a tweet <a href=\"https://chrisburnell.com/bowhead/\">https://chrisburnell.com/bowhead/</a>"
	},
	"target": "https://chrisburnell.com/bowhead/"
}
*/
type WebmentionIOFile struct {
	Links []WebmentionIOMention `json:"links"`
}

type WebmentionIOMention struct {
	Source       string               `json:"source"`
	Target       string               `json:"target"`
	Verified     bool                 `json:"verified"`
	VerifiedDate string               `json:"verified_date"`
	ID           int                  `json:"id"`
	Private      bool                 `json:"private"`
	Data         WebmentionIOData     `json:"data"`
	Activity     WebmentionIOActivity `json:"activity"`
}

type WebmentionIOActivity struct {
	Type         string `json:"type"`
	Sentence     string `json:"sentence"`
	SentenceHtml string `json:"sentence_html"`
}

type WebmentionIOData struct {
	Author      WebmentionIOAuthor `json:"author"`
	Url         string             `json:"url"`
	Name        string             `json:"name"`
	Content     string             `json:"content"`
	Published   string             `json:"published"`
	PublishedTs int                `json:"published_ts"`
}

type WebmentionIOAuthor struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Photo string `json:"photo"`
}

type WebmentionIOImporter struct {
}

func (wmio *WebmentionIOImporter) TryImport(data []byte) ([]*mf.IndiewebData, error) {
	var mentions WebmentionIOFile
	err := json.Unmarshal(data, &mentions)
	if err != nil {
		return nil, err
	}

	var converted []*mf.IndiewebData
	for _, wmiomention := range mentions.Links {
		converted = append(converted, convert(wmiomention))
	}

	return converted, nil
}

func convert(wmio WebmentionIOMention) *mf.IndiewebData {
	iType := typeOf(wmio)
	return &mf.IndiewebData{
		Author: mf.IndiewebAuthor{
			Name:    wmio.Data.Author.Name,
			Picture: wmio.Data.Author.Photo,
		},
		Name:         nameOf(wmio, iType),
		Content:      contentOf(wmio, iType),
		Published:    publishedDate(wmio),
		Url:          wmio.Data.Url,
		Source:       sourceOf(wmio),
		Target:       wmio.Target,
		IndiewebType: iType,
	}
}

// sourceOf returns wmio.Source unless it detects a silo link such as bridgy.
// In that case, it returns the data URL. This isn't entirely correct, as it technically never was the sender.
func sourceOf(wmio WebmentionIOMention) string {
	srcDomain := rest.Domain(wmio.Source)
	if common.Includes(rest.SiloDomains, srcDomain) {
		return wmio.Data.Url
	}

	return wmio.Source
}

func nameOf(wmio WebmentionIOMention, iType mf.MfType) string {
	if (iType == mf.TypeReply || iType == mf.TypeLike) && wmio.Data.Name == "" {
		return wmio.Data.Content
	}
	return wmio.Data.Name
}

func contentOf(wmio WebmentionIOMention, iType mf.MfType) string {
	content := wmio.Data.Content
	if iType == mf.TypeReply || (iType == mf.TypeLike && content == "") {
		content = wmio.Activity.Sentence
	}
	return common.Shorten(content)
}

// typeOf returns the mf.MfType from a wmio mention.
func typeOf(wmio WebmentionIOMention) mf.MfType {
	if wmio.Activity.Type == "like" {
		return mf.TypeLike
	}
	if wmio.Activity.Type == "bookmark" {
		return mf.TypeBookmark
	}
	if wmio.Activity.Type == "reply" {
		return mf.TypeReply
	}
	if wmio.Activity.Type == "link" {
		return mf.TypeLink
	}
	return mf.TypeMention
}

// publishedDate pries out Published or VerifiedDate from wmio with expecrted format: 2022-05-16T09:24:40+01:00
// This is the same as target format, but validate nonetheless
func publishedDate(wmio WebmentionIOMention) string {
	published := wmio.Data.Published
	if published == "" {
		published = wmio.VerifiedDate
	}

	return common.ToTime(published, mf.DateFormatWithTimeZone).Format(mf.DateFormatWithTimeZone)
}
