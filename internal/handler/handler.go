package handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/ihyaulhaq/url-shotener-BE/internal/middleware"
	"github.com/ihyaulhaq/url-shotener-BE/internal/service"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	userService UserServicer
	urlService  UrlService
	baseURL     string
	jwtSecret   string
}

type UserServicer interface {
	Login(ctx context.Context, email, password string) (service.LoginReturn, error)
	Register(ctx context.Context, username, email, password string) (service.LoginReturn, error)
	Logout(ctx context.Context, refreshToken string) error
}

type UrlService interface {
	CreateShortUrl(ctx context.Context, originalUrl string, userID *uuid.UUID) (service.ShortUrl, error)
	GetOriginalUrl(ctx context.Context, urlCode string) (service.ShortUrl, error)
	GetUserUrls(ctx context.Context, userID uuid.UUID) ([]service.ShortUrl, error)
}

func New(
	urlService UrlService,
	userService UserServicer,
	baseURL string,
	jwtSecret string,
) *Handler {
	return &Handler{
		urlService:  urlService,
		userService: userService,
		baseURL:     baseURL,
		jwtSecret:   jwtSecret,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.health)
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	// Optional auth — URL is linked to user if logged in, otherwise anonymous
	optionalAuth := middleware.OptionalAuth(h.jwtSecret)
	mux.Handle("POST /api/urls/shorten", optionalAuth(http.HandlerFunc(h.handleShorteningUrl)))

	// Required auth — only authenticated users can list their URLs
	requireAuth := middleware.RequireAuth(h.jwtSecret)
	mux.Handle("GET /api/urls", requireAuth(http.HandlerFunc(h.handleGetUserUrls)))

	mux.HandleFunc("GET /{shortUrl}", h.handleRedirectUrl)

	mux.HandleFunc("POST /login", h.handleUserLogin)
	mux.HandleFunc("POST /signup", h.handleUserSignUp)
	
	// requireAuth for logout to ensure a valid session is terminating (optional, but good practice)
	// We'll just let anyone hit /logout with a valid refresh_token, or require auth. 
	// The problem is if the access token is expired, they still need to logout. 
	// So usually POST /logout takes the refresh token without needing a valid access token.
	mux.HandleFunc("POST /logout", h.handleUserLogout)

	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
