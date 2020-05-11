package main

import (
	"log"
	"net/http"

	"github.com/grandima/user-count/handlers"
	"github.com/grandima/user-count/session"
	cache "github.com/grandima/user-count/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {

	sm := &session.SessionManager{}

	var redis = cache.NewRedisCache("redis://localhost:6379")

	defer redis.Client.Close()

	handler := handlers.Handler{Storage: redis, SessionManager: sm}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/go", handler.Handle)

	log.Fatal(http.ListenAndServe(":8003", r))
}
