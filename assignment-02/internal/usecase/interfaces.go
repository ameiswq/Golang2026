package usecase

import "assignment-02/pkg/modules"

type UserUsecase interface {
	GetUsers() ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user *modules.User) (int, error)
	UpdateUser(id int, user *modules.User) error
	DeleteUser(id int) error
}