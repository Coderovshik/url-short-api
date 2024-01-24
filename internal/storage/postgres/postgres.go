package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Coderovshik/url-short/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	conn *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*Storage, error) {
	const op = "storage.postgres.New"

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) Close() {
	s.conn.Close()
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) error {
	const op = "storage.postgres.SaveURL"

	query := "INSERT INTO url (url, alias) VALUES (@url, @alias)"
	args := pgx.NamedArgs{
		"url":   urlToSave,
		"alias": alias,
	}
	_, err := s.conn.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ExistAlias(ctx context.Context, alias string) bool {
	// const op = "storage.postgres.ExistAlias"

	query := "SELECT EXISTS (SELECT 1 FROM url WHERE alias = @alias) AS isExist"
	args := pgx.NamedArgs{
		"alias": alias,
	}

	row := s.conn.QueryRow(ctx, query, args)
	var isExist string
	err := row.Scan(&isExist)
	if err != nil {
		return false
	}

	if isExist == "true" {
		return true
	}

	return false
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	query := "SELECT url FROM url WHERE alias = @alias"
	args := pgx.NamedArgs{
		"alias": alias,
	}

	row := s.conn.QueryRow(ctx, query, args)
	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.postgres.DeleteURL"

	query := "DELETE FROM url WHERE alias = @alias"
	args := pgx.NamedArgs{
		"alias": alias,
	}

	_, err := s.conn.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "02000" {
				return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
