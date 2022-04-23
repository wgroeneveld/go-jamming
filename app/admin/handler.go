package admin

import (
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleGet(repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := mux.Vars(r)["domain"]
		rest.Json(w, repo.GetAllToModerate(domain))
	}
}

// TODO unit tests
// HandleApprove approves the Mention (by key in URL) and adds to the whitelist.
// Returns 200 OK with approved source/target or 404 if key is invalid.
func HandleApprove(c *common.Config, repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]

		approved := repo.Approve(key)
		if approved == nil {
			http.NotFound(w, r)
			return
		}

		c.AddToWhitelist(approved.AsMention().SourceDomain())
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Approved: %s", approved.AsMention().String())
	}
}

// HandleReject rejects the Mention (by key in URL) and adds to the blacklist.
// Returns 200 OK with rejected source/target or 404 if key is invalid.
func HandleReject(c *common.Config, repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]

		rejected := repo.Reject(key)
		if rejected == nil {
			http.NotFound(w, r)
			return
		}

		c.AddToBlacklist(rejected.AsMention().SourceDomain())
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Rejected: %s", rejected.AsMention().String())
	}
}
