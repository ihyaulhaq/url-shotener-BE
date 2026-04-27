package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"

	_ "github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

// handleShorteningUrl godoc
// @Summary      Shorten a URL
// @Description  Accepts a long URL and returns a shortened version
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        body body CreateURLRequest true "URL to shorten"
// @Success      201 {object} CreateURLResponse
// @Failure      400 {object} map[string]string "invalid request payload / validation error"
// @Failure      500 {object} map[string]string "internal server error"
// @Router       /api/urls/shorten [post]
func (h *Handler) handleShorteningUrl(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		OriginalUrl string `json:"original_url" validate:"required,url"`
	}

	type CreateURLResponse struct {
		ShortUrl    string `json:"short_url"`
		OriginalUrl string `json:"original_url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.ResponseWithError(w, 400, "invalid request payload")
		return
	}

	// validate struct tags
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.urlService.CreateShortUrl(r.Context(), params.OriginalUrl)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
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
		utils.ResponseWithError(w, http.StatusBadRequest, "url code is required")
		return
	}

	result, err := h.urlService.GetOriginalUrl(r.Context(), urlCode)
	if err != nil {
		utils.ResponseWithError(w, http.StatusNotFound, "short url not found")
		return
	}

	http.Redirect(w, r, result.OriginalUrl, http.StatusFound)

}
