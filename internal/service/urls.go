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

func (s *UrlService) CreateShortUrl(ctx context.Context, originalUrl string) (ShortUrl, error) {

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
