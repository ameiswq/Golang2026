package users

import (
	"errors"
	"fmt"
 	"time"
	"assignment-03/internal/repository/_postgres"
	"assignment-03/pkg/modules"
)

type Repository struct {
	db *_postgres.Dialect
	ExecTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{db: db, ExecTimeout: time.Second * 5,}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
  var users []modules.User
  err := r.db.DB.Select(&users, "select id, name, email, age, created_at from users;")
  if err != nil { return nil, err }
  fmt.Println(users)
  return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user, "select id, name, email, age, created_at from users where id=$1", id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
func (r *Repository) CreateUser(user *modules.User) (int, error) {
	var id int
	query := `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.DB.QueryRow(query, user.Name, user.Email, user.Age,).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateUser(id int, user *modules.User) error {
	result, err := r.db.DB.Exec("UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4",user.Name, user.Email, user.Age, id,)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}
func (r *Repository) DeleteUser(id int) error {
	result, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}