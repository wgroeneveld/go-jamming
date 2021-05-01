package pictures

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/db"
	_ "embed"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

//go:embed anonymous.jpg
var anonymous []byte

func init() {
	if anonymous == nil {
		log.Fatal().Msg("embedded anonymous image missing?")
	}
}

const (
	bridgy = "brid.gy"
)

// Handle handles picture GET calls.
// It does not validate the picture query as it's part of a composite key anyway.
func Handle(repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		picDomain := mux.Vars(r)["picture"]
		if picDomain == mf.Anonymous || picDomain == bridgy {
			servePicture(w, anonymous)
			return
		}

		picData := repo.GetPicture(picDomain)
		if picData == nil {
			http.NotFound(w, r)
			return
		}
		servePicture(w, picData)
	}
}

// servePicture writes an OK and raw bytes.
// For some reason, headers - although they should be there - aren't needed?
func servePicture(w http.ResponseWriter, bytes []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
