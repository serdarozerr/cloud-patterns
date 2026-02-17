package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func HttpLogger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "header", r.Header)
		start := time.Now()
		next(w, r)
		slog.Info("response", "header", w.Header(), "duration", time.Since(start))
	})
}
