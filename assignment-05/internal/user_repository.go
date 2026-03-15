package internal

import (
	"database/sql"
	"fmt"
	"strings"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetPaginatedUsers(filter UserFilter) (PaginatedResponse, error) {
	var (
		args []interface{}
		whereParts []string
		users []User
		totalCount int
		argPosition = 1
	)

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 5
	}
	if filter.OrderBy == "" {
		filter.OrderBy = "id"
	}
	if filter.OrderDir == "" {
		filter.OrderDir = "asc"
	}

	allowedOrderBy := map[string]string{
		"id": "id",
		"name": "name",
		"email": "email",
		"gender": "gender",
		"birth_date": "birth_date",
	}

	orderColumn, ok := allowedOrderBy[filter.OrderBy]
	if !ok {
		orderColumn = "id"
	}

	orderDir := strings.ToUpper(filter.OrderDir)
	if orderDir != "ASC" && orderDir != "DESC" {
		orderDir = "ASC"
	}

	if filter.ID != nil {
		whereParts = append(whereParts, fmt.Sprintf("id = $%d", argPosition))
		args = append(args, *filter.ID)
		argPosition++
	}

	if filter.Name != "" {
		whereParts = append(whereParts, fmt.Sprintf("name ILIKE $%d", argPosition))
		args = append(args, "%"+filter.Name+"%")
		argPosition++
	}

	if filter.Email != "" {
		whereParts = append(whereParts, fmt.Sprintf("email ILIKE $%d", argPosition))
		args = append(args, "%"+filter.Email+"%")
		argPosition++
	}

	if filter.Gender != "" {
		whereParts = append(whereParts, fmt.Sprintf("gender ILIKE $%d", argPosition))
		args = append(args, "%"+filter.Gender+"%")
		argPosition++
	}

	if filter.BirthDate != "" {
		whereParts = append(whereParts, fmt.Sprintf("birth_date = $%d", argPosition))
		args = append(args, filter.BirthDate)
		argPosition++
	}

	baseWhere := ""
	if len(whereParts) > 0 {
		baseWhere = " WHERE " + strings.Join(whereParts, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM users" + baseWhere
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return PaginatedResponse{}, err
	}

	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT id, name, email, gender, birth_date
		FROM users
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, baseWhere, orderColumn, orderDir, argPosition, argPosition+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return PaginatedResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return PaginatedResponse{}, err
	}

	return PaginatedResponse{
		Data: users,
		TotalCount: totalCount,
		Page: filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (r *UserRepository) GetCommonFriends(user1ID, user2ID int) ([]User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM user_friends uf1
		JOIN user_friends uf2 ON uf1.friend_id = uf2.friend_id
		JOIN users u ON u.id = uf1.friend_id
		WHERE uf1.user_id = $1
		  AND uf2.user_id = $2
		ORDER BY u.id
	`

	rows, err := r.db.Query(query, user1ID, user2ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}