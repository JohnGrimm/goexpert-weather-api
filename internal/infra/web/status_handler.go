package web

import (
	"encoding/json"
	"net/http"
)

type WebStatusHandler struct{}

func NewWebStatusHandler() *WebStatusHandler {
	return &WebStatusHandler{}
}

type Status struct {
	Status string `json:"status"`
}

func (h *WebStatusHandler) Get(w http.ResponseWriter, r *http.Request) {
	status := Status{
		Status: "UP",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
