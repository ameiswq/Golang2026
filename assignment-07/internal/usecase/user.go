package usecase

import (
	"errors"
	"assignment-07/internal/entity"
	"assignment-07/internal/usecase/repo"
	"assignment-07/utils"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, error) {
	exists, err := u.repo.ExistsByUsernameOrEmail(user.Username, user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	return u.repo.RegisterUser(user)
}

func (u *UserUseCase) LoginUser(input *entity.LoginUserDTO) (string, error) {
	userFromDB, err := u.repo.GetByUsername(input.Username)
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(userFromDB.Password, input.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(userFromDB.ID, userFromDB.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUseCase) GetMe(userID string) (*entity.User, error) {
	return u.repo.GetByID(userID)
}

func (u *UserUseCase) PromoteUser(id string) error {
	return u.repo.PromoteUser(id)
}