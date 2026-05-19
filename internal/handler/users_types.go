package handler

type userLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type userLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userLogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}
