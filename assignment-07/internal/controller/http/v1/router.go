package v1

import (
	"time"
	"assignment-07/internal/usecase"
	"assignment-07/utils"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, userUC usecase.UserInterface) {
	rl := utils.NewRateLimiter(5, time.Minute)

	api := handler.Group("/v1")
	{
		newUserRoutes(api, userUC, rl)
	}
}