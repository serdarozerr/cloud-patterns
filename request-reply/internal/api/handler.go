package api

import (
	"net/http"

	"github.com/serdarozerr/request-reply/internal/service/queue"
	m "github.com/serdarozerr/request-reply/pkg/middleware"
)

// var (
// 	authMiddleware = m.AuthMiddleware(config.NewConfig())
// )

func addUserRoutes(mux *http.ServeMux, producer *queue.Producer) {
	mux.HandleFunc("/api/v1/users", m.HttpLogger(users(producer)))
}



func NewRouter(producer *queue.Producer) http.Handler {
	mux := http.NewServeMux()
	addUserRoutes(mux,producer)

	return mux
}
