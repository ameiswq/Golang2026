package app

import (
	"context"
	"log"
	"net/http"
	"time"
	"assignment-03/internal/repository"
	"assignment-03/internal/repository/_postgres"
	"assignment-03/internal/usecase"
	"assignment-03/pkg/modules"
	"os"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbConfig := initPostgreConfig()
	pg := _postgres.NewPGDialect(ctx, dbConfig)
	repos := repository.NewRepositories(pg)
	userUC := usecase.NewUserUsecase(repos.UserRepository)
	router := newRouter(userUC)
	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func initPostgreConfig() *modules.PostgreConfig {
	host := os.Getenv("DB_HOST")
	if host == "" { host = "localhost" }
	port := os.Getenv("DB_PORT")
	if port == "" { port = "5432" }
	user := os.Getenv("DB_USER")
	if user == "" { user = "postgres" }
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" { pass = "postgres" }
	db := os.Getenv("DB_NAME")
	if db == "" { db = "godb" }
	ssl := os.Getenv("DB_SSLMODE")
	if ssl == "" { ssl = "disable" }
	return &modules.PostgreConfig{
		Host: host,
		Port: port,
		Username: user,
		Password: pass,
		DBName: db,
		SSLMode: ssl,
		ExecTimeout: 5 * time.Second,
	}
}