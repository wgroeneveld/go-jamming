
package index

import (
	"net/http"
	"fmt"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Printf("testje")
}

//func (s *server) handleIndex() http.HandlerFunc {

//}
