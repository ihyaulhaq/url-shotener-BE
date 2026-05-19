package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/ihyaulhaq/url-shotener-BE/internal/middleware"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

// handleShorteningUrl godoc
// @Summary      Shorten a URL
// @Description  Accepts a long URL and returns a shortened version. If authenticated, links the URL to the user.
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        body body CreateURLRequest true "URL to shorten"
// @Success      201 {object} CreateURLResponse
// @Failure      400 {object} map[string]string "invalid request payload / validation error"
// @Failure      500 {object} map[string]string "internal server error"
// @Router       /api/urls/shorten [post]
func (h *Handler) handleShorteningUrl(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := CreateURLRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.ErrInvalidRequest(err).Write(w)
		return
	}

	// validate struct tags
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		utils.ErrBadRequest("original url is required and must be a valid url").Write(w)
		return
	}

	//TODO:make the url_users works
	// Extract optional user ID from context (set by OptionalAuth middleware)
	var userIDPtr *uuid.UUID
	if uid, ok := middleware.UserIDFromContext(r.Context()); ok {
		userIDPtr = &uid
	}

	result, err := h.urlService.CreateShortUrl(r.Context(), params.OriginalUrl, userIDPtr)
	if err != nil {
		utils.ErrInternal(err).Write(w)
		return
	}

	utils.ResponseWithJSON(w, http.StatusCreated, CreateURLResponse{
		ShortUrl:    fmt.Sprintf("%s/%s", h.baseURL, result.UrlCode),
		OriginalUrl: result.OriginalUrl,
	})
}

// handleRedirectUrl godoc
// @Summary      Redirect to original URL
// @Description  Looks up the short code and redirects the client to the original URL
// @Tags         urls
// @Failure      400 {object} map[string]string "url code is required"
// @Failure      404 {object} map[string]string "short url not found"
// @Router       /{shortUrl} [get]
func (h *Handler) handleRedirectUrl(w http.ResponseWriter, r *http.Request) {
	urlCode := r.PathValue("shortUrl")
	if urlCode == "" {
		utils.ErrBadRequest("url code is required").Write(w)

		return
	}

	result, err := h.urlService.GetOriginalUrl(r.Context(), urlCode)
	if err != nil {
		utils.ErrNotFound("Short url not found").Write(w)
		return
	}

	http.Redirect(w, r, result.OriginalUrl, http.StatusFound)

}

// handleGetUserUrls godoc
// @Summary      Get user's URLs
// @Description  Returns all shortened URLs belonging to the authenticated user
// @Tags         urls
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} CreateURLResponse
// @Failure      401 {object} map[string]string "unauthorized"
// @Failure      500 {object} map[string]string "internal server error"
// @Router       /api/urls [get]
func (h *Handler) handleGetUserUrls(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		utils.ErrUnauthorized("unauthorized").Write(w)
		return
	}

	urls, err := h.urlService.GetUserUrls(r.Context(), userID)
	if err != nil {
		utils.ErrInternal(err).Write(w)
		return
	}

	type urlItem struct {
		ID          string `json:"id"`
		ShortUrl    string `json:"short_url"`
		OriginalUrl string `json:"original_url"`
		ClickCount  int32  `json:"click_count"`
		CreatedAt   string `json:"created_at"`
	}

	items := make([]urlItem, len(urls))
	for i, u := range urls {
		items[i] = urlItem{
			ID:          u.ID.String(),
			ShortUrl:    fmt.Sprintf("%s/%s", h.baseURL, u.UrlCode),
			OriginalUrl: u.OriginalUrl,
			ClickCount:  u.ClickCount,
			CreatedAt:   u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	utils.ResponseWithJSON(w, http.StatusOK, items)
}
