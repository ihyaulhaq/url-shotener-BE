package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ihyaulhaq/url-shotener-BE/internal/service"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

type UserLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) hanldleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if params.Email == "" || params.Password == "" {
		utils.ResponseWithError(w, http.StatusBadRequest, "email and password required")
		return
	}

	result, err := h.userService.Login(r.Context(), params.Email, params.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			utils.ResponseWithError(w, 404, "user not found")
		case errors.Is(err, service.ErrInvalidCredentials):
			utils.ResponseWithError(w, 401, "invalid credentials")
		default:
			utils.ResponseWithError(w, 500, "something went wrong")
		}
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, UserLoginResponse{
		Token:        result.AccessToken,
		RefreshToken: result.RefreshToken,
	})

}
