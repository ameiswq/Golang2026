package internal
import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgres() (*sql.DB, error) {
	host := "localhost"
	port := "5432"
	user := "postgres"
	password := "postgres"
	dbName := "practice5"
	sslmode := "disable"

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}