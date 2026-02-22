package usecase

import (
	"assignment-02/internal/repository"
	"assignment-02/pkg/modules"
)

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{repo: r}
}

func (u *userUsecase) GetUsers() ([]modules.User, error) {
	return u.repo.GetUsers()
}

func (u *userUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *userUsecase) CreateUser(user *modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

func (u *userUsecase) UpdateUser(id int, user *modules.User) error {
	return u.repo.UpdateUser(id, user)
}

func (u *userUsecase) DeleteUser(id int) error {
	return u.repo.DeleteUser(id)
}