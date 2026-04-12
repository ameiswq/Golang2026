package v1

import (
	"net/http"
	"assignment-07/internal/entity"
	"assignment-07/internal/usecase"
	"assignment-07/utils"
	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	u usecase.UserInterface
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.UserInterface, rl *utils.RateLimiter) {
	r := &userRoutes{u: u}
	h := handler.Group("/users")
	anonymousLimited := h.Group("/")
	anonymousLimited.Use(rl.Middleware())
	{
		anonymousLimited.POST("/", r.RegisterUser)
		anonymousLimited.POST("/login", r.LoginUser)
	}
	protected := h.Group("/")
	protected.Use(utils.JWTAuthMiddleware(), rl.Middleware())
	{
		protected.GET("/me", r.GetMe)
		protected.GET("/protected/hello", r.ProtectedHello)
	}
	admin := h.Group("/")
	admin.Use(utils.JWTAuthMiddleware(), rl.Middleware(), utils.RoleMiddleware("admin"))
	{
		admin.PATCH("/promote/:id", r.PromoteUser)
	}
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var dto entity.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}
	user := &entity.User{
		Username: dto.Username,
		Email: dto.Email,
		Password: hashedPassword,
		Role: "user",
		Verified: false,
	}
	createdUser, err := r.u.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": createdUser,
	})
}

func (r *userRoutes) LoginUser(c *gin.Context) {
	var dto entity.LoginUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.u.LoginUser(&dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	user, err := r.u.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": user.ID,
		"username": user.Username,
		"email": user.Email,
		"role": user.Role,
		"verified": user.Verified,
	})
}

func (r *userRoutes) PromoteUser(c *gin.Context) {
	id := c.Param("id")
	if err := r.u.PromoteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}

func (r *userRoutes) ProtectedHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}