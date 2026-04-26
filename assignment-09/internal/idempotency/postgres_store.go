package idempotency

import (
	"context"
	"database/sql"
	"errors"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{DB: db}
}

func (s *PostgresStore) StartProcessing(ctx context.Context, key string) (bool, error) {
	result, err := s.DB.ExecContext(ctx,
		`INSERT INTO idempotency_keys (key, status) VALUES ($1, 'processing') ON CONFLICT (key) DO NOTHING`,
		key,
	)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows == 1, nil
}

func (s *PostgresStore) Get(ctx context.Context, key string) (*CachedResponse, bool, error) {
	var status string
	var code sql.NullInt64
	var body sql.NullString

	err := s.DB.QueryRowContext(ctx,
		`SELECT status, response_code, response_body FROM idempotency_keys WHERE key = $1`,
		key,
	).Scan(&status, &code, &body)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	response := &CachedResponse{
		StatusCode: int(code.Int64),
		Body:       body.String,
		Completed:  status == "completed",
	}

	return response, true, nil
}

func (s *PostgresStore) Finish(ctx context.Context, key string, statusCode int, body string) error {
	_, err := s.DB.ExecContext(ctx,
		`UPDATE idempotency_keys SET status = 'completed', response_code = $2, response_body = $3 WHERE key = $1`,
		key,
		statusCode,
		body,
	)
	return err
}
