package pictures

import (
	"brainbaking.com/go-jamming/db"
	"github.com/gorilla/mux"
	"net/http"
)

// Handle handles picture GET calls.
// It does not validate the picture query as it's part of a composite key anyway.
func Handle(repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		picDomain := mux.Vars(r)["picture"]
		picData := repo.GetPicture(picDomain)
		if picData == nil {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		// TODO response headers? is this a jpeg, png, gif, webm? should we?
		w.Write(picData)
	}
}
