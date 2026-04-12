package usecase

import "assignment-07/internal/entity"

type UserInterface interface {
	RegisterUser(user *entity.User) (*entity.User, error)
	LoginUser(input *entity.LoginUserDTO) (string, error)
	GetMe(userID string) (*entity.User, error)
	PromoteUser(id string) error
}