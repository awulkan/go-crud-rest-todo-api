package main

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// Database holds a list of todos.
type Database struct {
	Todos []Todo
	sync.Mutex
}

// Todo represents a todo item.
type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Msg   string `json:"msg"`
	Done  bool   `json:"done"`
}

// The available characters to generate a random ID from.
// Characters must be valid in a URL path for the API to work.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Pseudo-random seed for the ID generator.
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// CreateDatabase returns a new Database instance.
func CreateDatabase() *Database {
	return &Database{}
}

// Count returns the amount of Todo items.
func (db *Database) Count() int {
	db.Lock()
	defer db.Unlock()
	return len(db.Todos)
}

// GetTodos returns all Todo items.
func (db *Database) GetTodos() *[]Todo {
	db.Lock()
	defer db.Unlock()
	return &db.Todos
}

// GetTodo returns the Todo item with the same ID as provided.
func (db *Database) GetTodo(id string) (*Todo, error) {
	db.Lock()
	defer db.Unlock()
	for _, todo := range db.Todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return &Todo{}, errors.New("Todo not found")
}

// DeleteTodo deletes the Todo item with the same ID as provided.
func (db *Database) DeleteTodo(id string) error {
	db.Lock()
	defer db.Unlock()
	for i, todo := range db.Todos {
		if todo.ID == id {
			// The most efficient way to splice slices I could find.
			copy(db.Todos[i:], db.Todos[i+1:])
			db.Todos = db.Todos[:len(db.Todos)-1]
			return nil
		}
	}
	return errors.New("Todo not found")
}

// AddTodo gives the provided Todo item an ID and appends it
// to the Database.
func (db *Database) AddTodo(t *Todo) {
	db.Lock()
	defer db.Unlock()

	t.ID = NewTodoID()
	db.Todos = append(db.Todos, *t)
}

// UpdateTodo gives the provided Todo item an ID and appends it
// to the Database.
func (db *Database) UpdateTodo(t *Todo) error {
	db.Lock()
	defer db.Unlock()
	for i := range db.Todos {
		if db.Todos[i].ID == t.ID {
			db.Todos[i] = *t
			return nil
		}
	}
	return errors.New("Todo not found")
}

// Populate fills the Database with a list of Todo items.
// If the dabase already contains todo items, or if the population
// failed, it will return an error.
func (db *Database) Populate() error {
	if db.Count() > 0 {
		return errors.New("The database is already populated")
	}

	db.AddTodo(&Todo{ID: NewTodoID(), Title: "Water flowers", Msg: "They're really dry...", Done: false})
	db.AddTodo(&Todo{ID: NewTodoID(), Title: "Pay bills", Msg: "Better get it done.", Done: false})
	db.AddTodo(&Todo{ID: NewTodoID(), Title: "Buy food", Msg: "Out of pasta.", Done: false})
	db.AddTodo(&Todo{ID: NewTodoID(), Title: "Make someone's day better", Msg: "I'm starting with the man in the mirror.", Done: false})

	if db.Count() != 5 {
		return errors.New("Populating the database failed")
	}
	return nil
}

// NewTodoID returns a simple pseudo-random string of letters and number.
// We don't need a too complicated generator for such a simple application as this.
func NewTodoID() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
