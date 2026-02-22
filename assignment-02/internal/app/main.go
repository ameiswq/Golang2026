package app

import (
	"context"
	"log"
	"net/http"
	"time"
	"assignment-02/internal/repository"
	"assignment-02/internal/repository/_postgres"
	"assignment-02/internal/usecase"
	"assignment-02/pkg/modules"
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
	return &modules.PostgreConfig{
		Host: "localhost",
		Port: "5432",
		Username: "postgres",
		Password: "postgres",
		DBName: "godb",
		SSLMode: "disable",
		ExecTimeout: 5 * time.Second,
	}
}