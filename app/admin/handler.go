package admin

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleGet(repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := mux.Vars(r)["domain"]
		rest.Json(w, repo.GetAllToModerate(domain))
	}
}

// TODO validate or not? see webmention.HandlePost
// TODO unit tests
func HandleApprove(c *common.Config, repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		wm := mf.Mention{
			Source: r.FormValue("source"),
			Target: r.FormValue("target"),
		}

		repo.Approve(wm)
		c.AddToWhitelist(wm.SourceDomain())
		w.WriteHeader(200)
	}
}

func HandleReject(c *common.Config, repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		wm := mf.Mention{
			Source: r.FormValue("source"),
			Target: r.FormValue("target"),
		}

		repo.Reject(wm)
		c.AddToBlacklist(wm.SourceDomain())
		w.WriteHeader(200)
	}
}
