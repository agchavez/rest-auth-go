package routers

import (
	"agchavez/go/rest-ws/handlers"
	"agchavez/go/rest-ws/server"
	"net/http"

	"github.com/gorilla/mux"
)

func AuthRouter(s server.Server, r *mux.Router) {
	r.HandleFunc("/auth/singin", handlers.SingUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
}
