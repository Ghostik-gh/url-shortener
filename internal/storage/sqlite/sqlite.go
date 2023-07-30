package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New" // for errors # op: Operation proc
	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = query.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlLong string, alias string) error {
	const op = "storage.sqlite.SaveURL"

	query, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = query.Exec(urlLong, alias)
	if err != nil {
		if err2, ok := err.(*sqlite.Error); ok && err2.Code() == 2067 {
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)

		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	query, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var url string

	err = query.QueryRow(alias).Scan(&url)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	query, err := s.db.Prepare("DELETE FROM url WHERE alias=?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	result, err := query.Exec(alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			return storage.ErrURLNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("%s: %v", op, "nothing to delete")
	}

	return nil
}
