package sqlite

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sql.DB
}

func New(storagePath string, dumpPath string) (*Storage, error) {
	const op = "sqlite.storage.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	b, err := os.ReadFile(dumpPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	schema := string(b)

	stmt, err := db.Prepare(schema)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		DB: db,
	}, nil
}
