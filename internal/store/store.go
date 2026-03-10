package store

import (
	"database/sql"

	"github.com/ihyaulhaq/url-shotener-BE/internal/database"
)

type Store struct {
	db *sql.DB
	*database.Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: database.New(db),
	}
}
