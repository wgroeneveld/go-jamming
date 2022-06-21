package rss

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"sort"
	"text/template"
	"time"
)

import _ "embed"

const (
	MaxRssItems = 50
)

//go:embed mentionsrss.xml
var mentionsrssTemplate []byte

type RssMentions struct {
	Domain string
	Date   time.Time
	Items  []*RssMentionItem
}

type RssMentionItem struct {
	ApproveURL string
	RejectURL  string
	Data       *mf.IndiewebData
}

func asTemplate(name string, data []byte) *template.Template {
	tmpl, err := template.New(name).Parse(string(data))
	if err != nil {
		log.Fatal().Err(err).Str("name", name).Msg("Template invalid")
	}
	return tmpl
}

func HandleGet(c *common.Config, repo db.MentionRepo) http.HandlerFunc {
	tmpl := asTemplate("mentionsRss", mentionsrssTemplate)

	return func(w http.ResponseWriter, r *http.Request) {
		domain := mux.Vars(r)["domain"]

		mentions := getLatestMentions(domain, repo, c)
		err := tmpl.Execute(w, RssMentions{
			Items:  mentions,
			Date:   time.Now(),
			Domain: domain,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("Unable to fill in dashboard template")
		}
	}
}

func getLatestMentions(domain string, repo db.MentionRepo, c *common.Config) []*RssMentionItem {
	toMod := repo.GetAllToModerate(domain).Data
	all := repo.GetAll(domain).Data

	var data []*RssMentionItem
	for _, v := range toMod {
		wm := v.AsMention()
		data = append(data, &RssMentionItem{
			Data:       v,
			ApproveURL: fmt.Sprintf("%sadmin/approve/%s/%s", c.BaseURL, c.Token, wm.Key()),
			RejectURL:  fmt.Sprintf("%sadmin/reject/%s/%s", c.BaseURL, c.Token, wm.Key()),
		})
	}
	for _, v := range all {
		data = append(data, &RssMentionItem{
			Data: v,
		})
	}

	// TODO this date is the published date, not the webmention received date!
	// This means it "might" disappear after the cutoff point in the RSS feed, and we don't store a received timestamp
	sort.Slice(data, func(i, j int) bool {
		return data[i].Data.PublishedDate().After(data[j].Data.PublishedDate())
	})
	if len(data) > MaxRssItems {
		return data[0:MaxRssItems]
	}
	return data
}
