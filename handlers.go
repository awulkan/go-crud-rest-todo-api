package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

// GetTodo returns one Todo item, which is determined by the provided ID.
func GetTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	if len(todoID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err := database.GetTodo(todoID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	todoJSON, err := json.Marshal(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(todoJSON)
}

// GetTodos returns all Todo items in the database.
func GetTodos(w http.ResponseWriter, r *http.Request) {
	todosJSON, err := json.Marshal(database.GetTodos())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(todosJSON)
}

// DeleteTodo removes the specified Todo item from the database.
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	if len(todoID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := database.DeleteTodo(todoID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// CreateTodo adds a new todo to the database. An ID is generated
// by the database and returned to the user.
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	t := Todo{}
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	database.AddTodo(&t)
	if len(t.ID) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todoID, err := json.Marshal(t.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(todoID)
}

// UpdateTodo updates an existing todo with the new values sent by the user.
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	if len(todoID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)

	t := Todo{
		ID: todoID,
	}
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = database.UpdateTodo(&t)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
