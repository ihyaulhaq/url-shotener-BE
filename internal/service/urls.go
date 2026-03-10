package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/ihyaulhaq/url-shotener-BE/internal/database"
	"github.com/ihyaulhaq/url-shotener-BE/internal/store"
)

type UrlService struct {
	store *store.Store
}

// constructor
func NewUrlService(s *store.Store) *UrlService {
	return &UrlService{store: s}
}

func (s *UrlService) CreateShortUrl(ctx context.Context, originalUrl string) (database.Url, error) {

	// validate url
	parsed, err := url.ParseRequestURI(originalUrl)
	if err != nil {
		return database.Url{}, fmt.Errorf("invaild url:%w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return database.Url{}, fmt.Errorf("url must use http or https")
	}

	// make the hash
	urlCode, err := s.generateUniqueCode(ctx, originalUrl)
	if err != nil {
		return database.Url{}, err
	}

	return s.store.CreateURL(ctx, database.CreateURLParams{
		UrlCode:     urlCode,
		OriginalUrl: originalUrl,
	})
}

func (s *UrlService) GetOriginalUrl(ctx context.Context, urlCode string) (database.Url, error) {
	if urlCode == "" {
		return database.Url{}, fmt.Errorf("url code is required")
	}

	result, err := s.store.GetURLByURLCode(ctx, urlCode)
	if err != nil {
		return database.Url{}, fmt.Errorf("short url not found")
	}

	updated, err := s.store.IncrementURLCount(context.Background(), result.ID)
	if err != nil {
		slog.Error("failed to increment click count", "urlCode", urlCode, "error", err)
		return result, nil
	}

	return updated, nil
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
