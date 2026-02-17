package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/serdarozerr/request-reply/internal/validators"
	v "github.com/serdarozerr/request-reply/pkg"
)

func users(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "request body is empty", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	data, errors, err := v.ValidateStruct[validators.CreateUser](r.Body)
	if err != nil {
		slog.Error("validation error", "err", err)
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		slog.Info("validation failed", "errors", errors)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errors,
		})
		return
	}

	slog.Info("user created", "data", data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "user created successfully",
		"data":    data,
	})
}
