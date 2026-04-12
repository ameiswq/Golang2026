package repo

import (
	"assignment-07/internal/entity"
	"assignment-07/pkg/postgres"
)

type UserRepo struct {
	PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{PG: pg}
}

func (r *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	if err := r.PG.Conn.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.PG.Conn.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByID(id string) (*entity.User, error) {
	var user entity.User
	if err := r.PG.Conn.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) PromoteUser(id string) error {
	return r.PG.Conn.Model(&entity.User{}).
		Where("id = ?", id).
		Update("role", "admin").Error
}

func (r *UserRepo) ExistsByUsernameOrEmail(username, email string) (bool, error) {
	var count int64
	err := r.PG.Conn.Model(&entity.User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error
	return count > 0, err
}