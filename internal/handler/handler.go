package handler

import (
	"context"
	"net/http"

	"github.com/ihyaulhaq/url-shotener-BE/internal/service"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

type Handler struct {
	userService UserServicer
	urlService  UrlService
	baseURL     string
}

type UserServicer interface {
	Login(ctx context.Context, email, password string) (service.LoginReturn, error)
	Register(ctx context.Context, username, email, password string) (service.LoginReturn, error)
}

type UrlService interface {
	CreateShortUrl(ctx context.Context, originalUrl string) (service.ShortUrl, error)
	GetOriginalUrl(ctx context.Context, urlCode string) (service.ShortUrl, error)
}

func New(
	urlService *service.UrlService,
	userService *service.UserService,
	baseURL string,
) *Handler {
	return &Handler{
		urlService:  urlService,
		userService: userService,
		baseURL:     baseURL,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.heatlh)

	mux.HandleFunc("POST /api/urls/shorten", h.handleShorteningUrl)
	mux.HandleFunc("GET /{shortUrl}", h.handleRedirectUrl)

	mux.HandleFunc("POST /login", h.hanldleUserLogin)
	mux.HandleFunc("POST /signup", h.hanldleUserSignUp)

	return mux
}

func (h *Handler) heatlh(w http.ResponseWriter, r *http.Request) {
	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})

}
