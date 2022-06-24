package main

import (
	"agchavez/go/rest-ws/routers"
	"agchavez/go/rest-ws/server"
	"context"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{Port: PORT, JWTSecret: JWT_SECRET, DatabaseURL: DATABASE_URL})

	if err != nil {
		log.Fatal("Error creating server: ", err)
	}
	s.StartServer(BindRoutes)

}

func BindRoutes(s server.Server, r *mux.Router) {
	routers.AuthRouter(s, r)
	routers.UserRouter(s, r)
}
