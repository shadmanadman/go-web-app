package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db := InitDB()

	defer db.Close()

	router := chi.NewRouter()

	router.Get("/users", GetUserHandler(db))
	router.Post("/users", CreateUserHandler(db))
	router.Get("/posts", GetPostHandler(db))
	router.Post("/posts", CreatePostHandler(db))

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", router)
}
