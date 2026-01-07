package user

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterNewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterNewUserResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
