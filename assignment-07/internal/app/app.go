package app

import (
	"os"
	v1 "assignment-07/internal/controller/http/v1"
	"assignment-07/internal/entity"
	"assignment-07/internal/usecase"
	"assignment-07/internal/usecase/repo"
	"assignment-07/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func Run() error {
	pg, err := postgres.New()
	if err != nil {
		return err
	}
	if err := pg.Conn.AutoMigrate(&entity.User{}); err != nil {
		return err
	}
	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)
	router := gin.Default()
	v1.NewRouter(router, userUseCase)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8090"
	}

	return router.Run(":" + port)
}