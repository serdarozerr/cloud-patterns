package api

import "net/http"

func tasks(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("tasks endpoint"))
}
