package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHitsRateLimitAfterSlammingRequests(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/endpoint", testFn).Methods("GET")
	r.Use(NewRateLimiter(5, 10).Middleware)
	ts := httptest.NewServer(r)

	t.Cleanup(ts.Close)
	statusCodes := []int{}

	for i := 0; i <= 10; i++ {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/endpoint", ts.URL), nil)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		statusCodes = append(statusCodes, resp.StatusCode)
	}
	assert.Contains(t, statusCodes, 429)
}

func TestDoesNotHitRateLimitOfSecondEndpointAfterSlammingFirst(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/endpoint1", testFn).Methods("GET")
	r.HandleFunc("/endpoint2", testFn).Methods("GET")
	r.Use(NewRateLimiter(5, 10).Middleware)
	ts := httptest.NewServer(r)

	t.Cleanup(ts.Close)

	for i := 0; i <= 10; i++ {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/endpoint1", ts.URL), nil)
		client.Do(req)
	}
	for i := 0; i <= 5; i++ {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/endpoint2", ts.URL), nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}
}

func testFn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
