package user

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrRefreshTokenNotFound  = errors.New("refresh token not found")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrRefreshTokenExpired   = errors.New("refresh token expired")
	ErrTokenNotValid         = errors.New("token not valid")
	ErrPasswordMismatch      = errors.New("password mismatch")
	ErrTokenRevoked          = errors.New("token revoked")
)
