package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

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

func (h *Handler) handleRedirectUrl(w http.ResponseWriter, r *http.Request) {
	urlCode := r.PathValue("urlCode") // or chi: chi.URLParam(r, "urlCode")
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
