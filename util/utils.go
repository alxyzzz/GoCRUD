package util

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func SendJson(w http.ResponseWriter, res Response, status int) {
	data, err := json.Marshal(res)
	if err != nil {
		slog.Error("error when doing marshal of json", "error", err)
		SendJson(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("error when sending response", "error", err)
		return
	}
}
