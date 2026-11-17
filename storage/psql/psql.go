package postgresql

import (
	"awesomeProject/storage"
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	Db *sql.DB
}

func New(ctx context.Context, ConnectionString string) (*Storage, error) {
	const op = "storage/psql/psql.go/new"

	db, err := sql.Open("pgx", ConnectionString)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS url
	(
	id SERIAL PRIMARY KEY,
	alias TEXT NOT NULL UNIQUE,
	url TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`)

	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	storage := Storage{
		Db: db,
	}

	return &storage, nil
}

func (s *Storage) SaveUrl(saveUrl string, alias string, ctx context.Context) (int64, error) {
	const op = "storage/psql/psql.go/SaveUrl"
	const delete_flag = false
	var id int64

	err := s.Db.QueryRowContext(ctx, `INSERT INTO url(alias,url,delete_flag) VALUES($1,$2,$3) RETURNING id`, alias, saveUrl, delete_flag).Scan(&id)

	if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return id, nil
}

func (s Storage) GetUrl(alias string, ctx context.Context) (string, error) {
	const op = "storage/psql/psql.go/GetUrl"

	var url string

	err := s.Db.QueryRowContext(ctx, `SELECT url FROM url WHERE alias = $1`, alias).Scan(url)

	if err != nil {
		return "", fmt.Errorf("%s Ошибка запроса: %s", op, err)
	}

	return url, nil
}

func (s Storage) DeleteUrl(ctx context.Context, alias string) error {
	const op = "storage/psql/psql.go/DeleteUrl"

	_, err := s.Db.ExecContext(ctx, `UPDATE url SET delete_flag = true WHERE alias=$1`, alias)

	if err != nil {
		return fmt.Errorf("%s ошибка удаления: %w", op, err)
	}
	return nil
}
