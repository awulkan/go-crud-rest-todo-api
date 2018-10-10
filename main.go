package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var database = CreateDatabase()

func main() {
	chi := chi.NewRouter()

	chi.Use(middleware.Logger)
	chi.Use(middleware.Recoverer)

	chi.Get("/todo", GetTodos)
	chi.Get("/todo/{id}", GetTodo)
	chi.Post("/todo", CreateTodo)
	chi.Put("/todo/{id}", UpdateTodo)
	chi.Delete("/todo/{id}", DeleteTodo)

	// Populate the todo list
	database.Populate()

	server := &http.Server{
		Addr:         ":3000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      chi,
	}

	log.Println("Starting server at: https://localhost" + server.Addr)
	log.Fatal(server.ListenAndServe())
}
