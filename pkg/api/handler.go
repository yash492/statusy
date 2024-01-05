package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Status  int  `json:"status"`
	Data    any  `json:"data"`
	IsError bool `json:"is_error,omitempty"`
}

type Handler func(w http.ResponseWriter, r *http.Request) *Response

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := fn(w, r)
	if res == nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)

	if res.Status == http.StatusNoContent {
		return
	}

	err := json.NewEncoder(w).Encode(res.Data)
	if err != nil {
		slog.Error(err.Error())
	}
}
