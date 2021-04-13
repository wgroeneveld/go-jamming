package app

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// mimic ResponseWriter's WriteHeader to capture the code
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Logs ip (for those pingback bad boys to put in jail), request url, method, and response status code.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logWriter := &loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(logWriter, r)
		log.Info().
			Str("url", r.RequestURI).
			Str("ip", ipFrom(r)).
			Str("method", r.Method).
			Int("status", logWriter.statusCode).
			Msg("handled")
	})
}
