package admin

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"text/template"
)

import _ "embed"

//go:embed dashboard.html
var dashboardTemplate []byte

//go:embed moderated.html
var moderatedTemplate []byte

type dashboardMention struct {
	Source     string
	Target     string
	Content    string
	ApproveURL string
	RejectURL  string
}

type domainMention struct {
	Name        string
	MentionsURL string
}

type dashboardData struct {
	Config   string
	Mentions map[domainMention][]dashboardMention
}

type dashboardModerated struct {
	Action      string
	Item        string
	RedirectURL string
}

func indiewebDataToDashboardMention(c *common.Config, dbMentions []*mf.IndiewebData) []dashboardMention {
	var mentions []dashboardMention
	for _, dbMention := range dbMentions {
		wm := dbMention.AsMention()
		// TODO move this to somewhere else? the wm? duplicate in notifier.go
		approveUrl := fmt.Sprintf("%sadmin/approve/%s/%s", c.BaseURL, c.Token, wm.Key())
		rejectUrl := fmt.Sprintf("%sadmin/reject/%s/%s", c.BaseURL, c.Token, wm.Key())

		mentions = append(mentions, dashboardMention{
			Source:     dbMention.Source,
			Target:     dbMention.Target,
			Content:    dbMention.Content,
			ApproveURL: approveUrl,
			RejectURL:  rejectUrl,
		})
	}

	return mentions
}

func getDashboardData(c *common.Config, repo db.MentionRepo) *dashboardData {
	data := &dashboardData{
		Config:   c.String(),
		Mentions: map[domainMention][]dashboardMention{},
	}
	for _, domain := range c.AllowedWebmentionSources {
		domainKey := domainMention{
			Name:        domain,
			MentionsURL: fmt.Sprintf("%swebmention/%s/%s", c.BaseURL, domain, c.Token),
		}
		data.Mentions[domainKey] = indiewebDataToDashboardMention(c, repo.GetAllToModerate(domain).Data)
	}

	return data
}

func asTemplate(name string, data []byte) *template.Template {
	tmpl, err := template.New(name).Parse(string(data))
	if err != nil {
		log.Fatal().Err(err).Str("name", name).Msg("Template invalid")
	}
	return tmpl
}

func HandleGet(c *common.Config, repo db.MentionRepo) http.HandlerFunc {
	tmpl := asTemplate("dashboard", dashboardTemplate)

	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, getDashboardData(c, repo))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("Unable to fill in dashboard template")
		}
	}
}

func HandleGetToApprove(repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := mux.Vars(r)["domain"]
		rest.Json(w, repo.GetAllToModerate(domain))
	}
}

// HandleApprove approves the Mention (by key in URL) and adds to the whitelist.
// Returns 200 OK with approved source/target or 404 if key is invalid.
func HandleApprove(c *common.Config, repo db.MentionRepo) http.HandlerFunc {
	tmpl := asTemplate("moderated", moderatedTemplate)

	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]

		approved := repo.Approve(key)
		if approved == nil {
			http.NotFound(w, r)
			return
		}

		c.AddToWhitelist(approved.AsMention().SourceDomain())
		err := tmpl.Execute(w, asDashboardModerated("Approved", approved, c))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("Unable to fill in dashboard template")
		}
	}
}

// HandleReject rejects the Mention (by key in URL) and adds to the blacklist.
// Returns 200 OK with rejected source/target or 404 if key is invalid.
func HandleReject(c *common.Config, repo db.MentionRepo) http.HandlerFunc {
	tmpl := asTemplate("moderated", moderatedTemplate)

	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]

		rejected := repo.Reject(key)
		if rejected == nil {
			http.NotFound(w, r)
			return
		}

		c.AddToBlacklist(rejected.AsMention().SourceDomain())
		err := tmpl.Execute(w, asDashboardModerated("Rejected", rejected, c))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("Unable to fill in dashboard template")
		}
	}
}

func asDashboardModerated(action string, mention *mf.IndiewebData, c *common.Config) dashboardModerated {
	return dashboardModerated{
		Action:      action,
		Item:        mention.AsMention().String(),
		RedirectURL: fmt.Sprintf("%sadmin/%s", c.BaseURL, c.Token),
	}
}
