package rss

import (
	"brainbaking.com/go-jamming/common"
	"encoding/xml"
	"errors"
	"github.com/rs/zerolog/log"
	"html/template"
	"time"
)

// someone already did this for me, yay! https://siongui.github.io/2015/03/03/go-parse-web-feed-rss-atom/
type Rss2 struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	// Required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	// Optional
	PubDate  string `xml:"channel>pubDate"`
	ItemList []Item `xml:"channel>item"`
}

type Item struct {
	// Required
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description template.HTML `xml:"description"`
	// Optional
	Content  template.HTML `xml:"encoded"`
	PubDate  string        `xml:"pubDate"`
	Comments string        `xml:"comments"`
}

func (itm Item) PubDateAsTime() time.Time {
	// format: Tue, 16 Mar 2021 17:07:14 +0000
	t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", itm.PubDate)
	if err != nil {
		log.Warn().Str("pubDate", itm.PubDate).Msg("Incorrectly formatted RSS date, reverting to now")
		return common.Now()
	}
	return t
}

type Link struct {
	Href string `xml:"href,attr"`
}

type Author struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type Entry struct {
	Title   string `xml:"title"`
	Summary string `xml:"summary"`
	Content string `xml:"content"`
	Id      string `xml:"id"`
	Updated string `xml:"updated"`
	Link    Link   `xml:"link"`
	Author  Author `xml:"author"`
}

func ParseFeed(content []byte) (Rss2, error) {
	v := Rss2{}
	err := xml.Unmarshal(content, &v)
	if err != nil {
		return v, err
	}

	if v.Version == "2.0" {
		for i, _ := range v.ItemList {
			if v.ItemList[i].Content != "" {
				v.ItemList[i].Description = v.ItemList[i].Content
			}
		}
		return v, nil
	}

	return v, errors.New("not RSS 2.0")
}
