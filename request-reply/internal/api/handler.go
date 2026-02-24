package api

import (
	"net/http"

	"github.com/serdarozerr/request-reply/internal/service"
	m "github.com/serdarozerr/request-reply/pkg/middleware"
)

// var (
// 	authMiddleware = m.AuthMiddleware(config.NewConfig())
// )

func addUserRoutes(mux *http.ServeMux, producer *service.Producer) {
	mux.HandleFunc("/api/v1/users", m.HttpLogger(users(producer)))
}



func NewRouter(producer *service.Producer) http.Handler {
	mux := http.NewServeMux()
	addUserRoutes(mux,producer)

	return mux
}
