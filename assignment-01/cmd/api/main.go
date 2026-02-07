package main

import (
	"log"
	"net/http"
	"github.com/ameiswq/Golang2026/assignment-01/internal/storage"
	"github.com/ameiswq/Golang2026/assignment-01/internal/handlers"
	"github.com/ameiswq/Golang2026/assignment-01/internal/middleware"
)

func main() {
	store := storage.NewStore()
	handler := handlers.NewTaskHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handler.HandleTasks)

	finalHandler := middleware.LoggingMiddleware(middleware.AuthMiddleware(mux))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		log.Fatal(err)
	}
}