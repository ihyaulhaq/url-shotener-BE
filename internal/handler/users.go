package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ihyaulhaq/url-shotener-BE/internal/service"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

// handleUserLogin godoc
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
func (h *Handler) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	params := userLoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.ErrInvalidRequest(err).Write(w)
		return
	}
	if params.Email == "" || params.Password == "" {
		utils.ErrBadRequest("email and password are required").Write(w)
		return
	}

	result, err := h.userService.Login(r.Context(), params.Email, params.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			utils.ErrNotFound("user").Write(w)
		case errors.Is(err, service.ErrInvalidCredentials):
			utils.ErrUnauthorized("invalid credentials").Write(w)
		default:
			utils.ErrInternal(err).Write(w)
		}
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, userLoginResponse{
		Token:        result.AccessToken,
		RefreshToken: result.RefreshToken,
	})

}

// handleUserSignUp godoc
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
func (h *Handler) handleUserSignUp(w http.ResponseWriter, r *http.Request) {

	params := createUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.ErrInvalidRequest(err).Write(w)
		return
	}
	if params.Email == "" || params.Password == "" {
		utils.ErrBadRequest("email and password required").Write(w)
		return
	}

	result, err := h.userService.Register(r.Context(), params.Username, params.Email, params.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailTaken):
			utils.ErrConflict("email already taken").Write(w)
		default:
			utils.ErrInternal(err).Write(w)
		}
		return
	}

	utils.ResponseWithJSON(w, http.StatusCreated, userLoginResponse{
		Token:        result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

// handleUserLogout godoc
// @Summary      Logout user
// @Description  Revokes a user's refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body userLogoutRequest true "Refresh token to revoke"
// @Success      200 {object} map[string]string "success message"
// @Failure      400 {object} utils.Envelope "invalid request payload"
// @Failure      500 {object} utils.Envelope "something went wrong"
// @Router       /logout [post]
func (h *Handler) handleUserLogout(w http.ResponseWriter, r *http.Request) {
	params := userLogoutRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.ErrInvalidRequest(err).Write(w)
		return
	}

	if params.RefreshToken == "" {
		utils.ErrBadRequest("refresh token is required").Write(w)
		return
	}

	err := h.userService.Logout(r.Context(), params.RefreshToken)
	if err != nil {
		utils.ErrInternal(err).Write(w)
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{
		"message": "successfully logged out",
	})
}
