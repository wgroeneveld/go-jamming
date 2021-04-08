
package rest

import (
	"net/http"
)

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "400 bad request", http.StatusBadRequest)
}

func Accept(w http.ResponseWriter) {
	w.WriteHeader(202)
	w.Write([]byte("Thanks, bro. Will send these webmentions soon, pinky swear!"))
}
