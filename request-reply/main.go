package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/serdarozerr/request-reply/internal/api"
)

func configureLogger() {

	handler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))

}

func main() {
	configureLogger()

	slog.Info("Starting server on port 8080")

	s := http.Server{
		Addr:    ":8080",
		Handler: api.NewRouter(),
	}

	if err := s.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "port", 8080)
	}

}
