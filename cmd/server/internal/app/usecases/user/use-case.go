package user

import (
	"context"
	"errors"
	"fmt"
	domain "server/internal/app/domain/user"
	hasher "server/internal/pkg/hash/argon2"
	"server/internal/pkg/token"
	"time"
)

type Repository interface {
	CreateNewUser(ctx context.Context, user *domain.User) (int64, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetTokens(ctx context.Context, userId int64) (*token.Tokens, error)
	AddTokens(ctx context.Context, userId int64, token *token.Tokens) error
	UpdateTokens(ctx context.Context, userId int64, token *token.Tokens) error
}

type User struct {
	repo Repository
}

func New(repo Repository) *User {
	return &User{repo: repo}
}

// RegisterNewUser Creates new user, JWT and Refresh token
func (u *User) RegisterNewUser(ctx context.Context, username, password string) (*token.Tokens, error) {
	// create new empty User object
	user := domain.NewUser()

	// hash user password
	hashedPassword, err := hasher.HashString(password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	user.Username = username

	// create new User in database
	userID, err := u.repo.CreateNewUser(ctx, user)
	if err != nil {
		if !errors.Is(err, domain.ErrUsernameAlreadyExists) {
			return nil, fmt.Errorf("failed to create new user: %w", err)
		}
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	// create tokens
	return u.createTokens(ctx, userID)
}

// Login compare hashed password from database with password from request, create tokens
func (u *User) Login(ctx context.Context, username, password string) (*token.Tokens, error) {
	user, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return nil, fmt.Errorf("failed to find user by username: %w", err)
		}
		return nil, fmt.Errorf("got unexpected err: %w", err)
	}
	ok, err := hasher.VerifyString(password, user.Password)
	if !ok {
		return nil, domain.ErrPasswordMismatch
	}
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}

	return u.createTokens(ctx, user.ID)
}

// Authenticate verify JWT token and get UserID from it
func (u *User) Authenticate(t string) (int64, error) {
	jwt, err := token.VerifyJWT(t)
	if err != nil {
		return 0, domain.ErrTokenNotValid
	}
	return jwt.UserID, nil
}

// RefreshJWTToken refresh JWT token by refresh token from request, create tokens
func (u *User) RefreshJWTToken(ctx context.Context, jwt, refreshToken string) (*token.Tokens, error) {
	// get UserId from jwt payload
	claim, err := token.VerifyJWT(jwt) // fixme
	if err != nil {
		return nil, err
	}

	// get tokens by UserID
	tokens, err := u.repo.GetTokens(ctx, claim.UserID)
	if err != nil {
		return nil, err
	}

	// check if tokens not revoked
	if tokens.Revoked {
		return nil, domain.ErrTokenRevoked
	}

	if time.Now().After(tokens.RefreshTokenExpAt) {
		return nil, domain.ErrRefreshTokenExpired
	}

	// verify hashed tokens
	ok, err := hasher.VerifyString(refreshToken, tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify refresh token: %w", err)
	}
	if !ok {
		return nil, domain.ErrInvalidRefreshToken
	}
	return u.updateTokens(ctx, claim.UserID)
}

func (u *User) createTokens(ctx context.Context, userID int64) (*token.Tokens, error) {
	// create new empty Tokens
	t := token.NewTokens(userID)

	// create jwt
	jwt, err := token.CreateNewJWT(userID)
	if err != nil {
		return nil, err
	}
	t.AddJWTToken(jwt)

	// create refresh token
	refreshToken := token.CreateRefreshToken()

	// add hashed refresh token in database
	t.AddRefreshToken(refreshToken)

	// hash refresh token
	hashedRefreshToken, err := hasher.HashString(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}
	t.RefreshToken = hashedRefreshToken

	// save t in database
	err = u.repo.AddTokens(ctx, userID, t)
	if err != nil {
		return nil, fmt.Errorf("failed to add t: %w", err)
	}

	t.RefreshToken = refreshToken
	return t, nil
}

func (u *User) updateTokens(ctx context.Context, userID int64) (*token.Tokens, error) {
	// create new empty Tokens
	tokens := token.NewTokens(userID)

	// create jwt
	jwt, err := token.CreateNewJWT(userID)
	if err != nil {
		return nil, err
	}
	tokens.AddJWTToken(jwt)

	// create refresh token
	refreshToken := token.CreateRefreshToken()

	// add hashed refresh token in database
	tokens.AddRefreshToken(refreshToken)

	// hash refresh token
	hashedRefreshToken, err := hasher.HashString(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}
	tokens.RefreshToken = hashedRefreshToken

	// save tokens in database
	err = u.repo.UpdateTokens(ctx, userID, tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to update  tokens: %w", err)
	}

	tokens.RefreshToken = refreshToken
	return tokens, nil
}
