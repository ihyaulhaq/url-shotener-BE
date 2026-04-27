package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ihyaulhaq/url-shotener-BE/internal/service"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"

	_ "github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

// hanldleUserLogin godoc
// @Summary      Login user
// @Description  Authenticates a user and returns an access token and refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body userLoginRequest true "User credentials"
// @Success      200 {object} userLoginResponse
// @Failure      400 {object} utils.Envelope "invalid request payload"
// @Failure      401 {object} utils.Envelope "invalid credentials"
// @Failure      404 {object} utils.Envelope "user not found"
// @Failure      500 {object} utils.Envelope "something went wrong"
// @Router       /login [post]
func (h *Handler) hanldleUserLogin(w http.ResponseWriter, r *http.Request) {
	params := userLoginRequest{}
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

	utils.ResponseWithJSON(w, http.StatusOK, userLoginResponse{
		Token:        result.AccessToken,
		RefreshToken: result.RefreshToken,
	})

}

// hanldleUserSignUp godoc
// @Summary      Register a new user
// @Description  Creates a new user account with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body createUserRequest true "User registration details"
// @Success      201 {object} userLoginResponse
// @Failure      400 {object} map[string]string "email and password required / invalid request payload"
// @Failure      409 {object} map[string]string "email already taken"
// @Failure      500 {object} map[string]string "something went wrong"
// @Router       /signup [post]
func (h *Handler) hanldleUserSignUp(w http.ResponseWriter, r *http.Request) {
	type response struct {
	}

	params := createUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if params.Email == "" || params.Password == "" {
		utils.ResponseWithError(w, http.StatusBadRequest, "email and password required")
		return
	}

	result, err := h.userService.Register(r.Context(), params.Username, params.Email, params.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailTaken):
			utils.ResponseWithError(w, http.StatusConflict, "email already taken")
		default:
			utils.ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}

	utils.ResponseWithJSON(w, http.StatusCreated, userLoginResponse{
		Token:        result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}
