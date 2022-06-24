package routers

import (
	"agchavez/go/rest-ws/handlers"
	"agchavez/go/rest-ws/server"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRouter(s server.Server, r *mux.Router) {
	r.HandleFunc("/user", handlers.ListUsersHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", handlers.GetUserHandler(s)).Methods(http.MethodGet)
}
