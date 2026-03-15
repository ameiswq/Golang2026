package main

import (
	"log"
	"net/http"
	"assignment-05/internal"
)

func main() {
	database, err := internal.NewPostgres()
	if err != nil {
		log.Fatal("db connection error:", err)
	}
	defer database.Close()

	userRepo := internal.NewUserRepository(database)
	userHandler := internal.NewUserHandler(userRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.GetUsers)
	mux.HandleFunc("/users/common-friends", userHandler.GetCommonFriends)

	log.Println("server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}