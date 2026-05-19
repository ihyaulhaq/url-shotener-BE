package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/ihyaulhaq/url-shotener-BE/internal/database"
	"github.com/ihyaulhaq/url-shotener-BE/internal/store"
)

type UrlService struct {
	store *store.Store
}

type ShortUrl struct {
	ID          uuid.UUID
	UrlCode     string
	OriginalUrl string
	ClickCount  int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// constructor
func NewUrlService(s *store.Store) *UrlService {
	return &UrlService{store: s}
}

func (s *UrlService) CreateShortUrl(ctx context.Context, originalUrl string, userID *uuid.UUID) (ShortUrl, error) {

	// validate url
	parsed, err := url.ParseRequestURI(originalUrl)
	if err != nil {
		return ShortUrl{}, fmt.Errorf("invaild url:%w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ShortUrl{}, fmt.Errorf("url must use http or https")
	}

	// make the hash
	urlCode, err := s.generateUniqueCode(ctx, originalUrl)
	if err != nil {
		return ShortUrl{}, err
	}

	dbUrl, err := s.store.CreateURL(ctx, database.CreateURLParams{
		UrlCode:     urlCode,
		OriginalUrl: originalUrl,
	})
	if err != nil {
		return ShortUrl{}, err
	}

	// If user is logged in, link the URL to the user
	if userID != nil {
		_, err := s.store.CreateURLUser(ctx, database.CreateURLUserParams{
			UrlID:  dbUrl.ID,
			UserID: *userID,
		})
		if err != nil {
			slog.Error("failed to link url to user", "urlID", dbUrl.ID, "userID", *userID, "error", err)
		}
	}

	return ShortUrl{
		ID:          dbUrl.ID,
		UrlCode:     dbUrl.UrlCode,
		OriginalUrl: dbUrl.OriginalUrl,
		ClickCount:  dbUrl.ClickCount,
		CreatedAt:   dbUrl.CreatedAt,
		UpdatedAt:   dbUrl.UpdatedAt,
	}, nil

}

func (s *UrlService) GetOriginalUrl(ctx context.Context, urlCode string) (ShortUrl, error) {
	if urlCode == "" {
		return ShortUrl{}, fmt.Errorf("url code is required")
	}

	result, err := s.store.GetURLByURLCode(ctx, urlCode)
	if err != nil {
		return ShortUrl{}, fmt.Errorf("short url not found")
	}
	updated, err := s.store.IncrementURLCount(context.Background(), result.ID)
	if err != nil {
		slog.Error("failed to increment click count", "urlCode", urlCode, "error", err)
		return ShortUrl{}, nil
	}

	url := ShortUrl{
		ID:          updated.ID,
		UrlCode:     updated.UrlCode,
		OriginalUrl: updated.OriginalUrl,
		ClickCount:  updated.ClickCount,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}

	return url, nil
}

func (s *UrlService) GetUserUrls(ctx context.Context, userID uuid.UUID) ([]ShortUrl, error) {
	dbUrls, err := s.store.GetURLsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user urls: %w", err)
	}

	urls := make([]ShortUrl, len(dbUrls))
	for i, u := range dbUrls {
		urls[i] = ShortUrl{
			ID:          u.ID,
			UrlCode:     u.UrlCode,
			OriginalUrl: u.OriginalUrl,
			ClickCount:  u.ClickCount,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
		}
	}

	return urls, nil
}

func (s *UrlService) generateUniqueCode(ctx context.Context, input string) (string, error) {
	const maxRetries = 5
	for i := range maxRetries {
		salted := fmt.Sprintf("%s:%d", input, i)
		hash := sha256.Sum256([]byte(salted))
		code := fmt.Sprintf("%x", hash[:6])

		if !s.IsCollision(ctx, code) {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique url code after %d attempts", maxRetries)
}

func (s *UrlService) IsCollision(ctx context.Context, str string) bool {
	_, err := s.store.GetURLByURLCode(ctx, str)
	return err == nil
}
