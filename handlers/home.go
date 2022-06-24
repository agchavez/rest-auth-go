package handlers

import (
	"agchavez/go/rest-ws/server"
	"encoding/json"
	"net/http"
)

type HomeResponse struct {
	Menssage string `json:"message"`
	Status   bool   `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(
			HomeResponse{
				Menssage: "Hello, world!",
				Status:   true,
			})
	}

}
