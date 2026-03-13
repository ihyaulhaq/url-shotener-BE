package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
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

// ++++++++++++++++++++++++++ AUTH +++++++++++++++++++++++++++++++++
func (s *Store) CreateRefreshTokenForUser(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) (database.RefreshToken, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return database.RefreshToken{}, err
	}
	defer tx.Rollback()

	qtx := s.Queries.WithTx(tx)

	// if you had to invalidate old tokens first, you'd do it here atomically
	// qtx.DeleteExistingRefreshTokens(ctx, userID)

	refreshToken, err := qtx.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return database.RefreshToken{}, err
	}

	return refreshToken, tx.Commit()
}
