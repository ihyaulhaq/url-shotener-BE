package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

func (h *Handler) hanldleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request payload")
	}

	if params.Email == "" || params.Password == "" {
		utils.ResponseWithError(w, http.StatusBadRequest, "email and password required")
	}
	result, err := h.userService.Login(params.Email, params.Password)

}
