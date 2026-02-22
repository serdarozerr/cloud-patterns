package api

import (
	"net/http"

	"github.com/serdarozerr/request-reply/internal/config"
	m "github.com/serdarozerr/request-reply/pkg/middleware"
)

var (
	authMiddleware = m.AuthMiddleware(config.NewConfig())
)

func addUserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/users", m.HttpLogger(users))
}



func NewRouter() http.Handler {
	mux := http.NewServeMux()
	addUserRoutes(mux)

	return mux
}
