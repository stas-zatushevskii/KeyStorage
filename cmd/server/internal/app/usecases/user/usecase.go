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
	user.Password = hasher.HashString(password)
	user.Username = username

	// create new User in database
	userID, err := u.repo.CreateNewUser(ctx, user)
	if err != nil {
		if !errors.Is(err, domain.ErrUsernameAlreadyExists) {
			return nil, fmt.Errorf("failed to create new user: %w", err)
		}
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
		return 0, err // fixme: require error type check
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
	if hasher.HashString(tokens.RefreshToken) != refreshToken {
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
	nonHashedRefreshToken := t.RefreshToken
	t.RefreshToken = hasher.HashString(t.RefreshToken)

	// save t in database
	err = u.repo.AddTokens(ctx, userID, t)
	if err != nil {
		return nil, fmt.Errorf("failed to add t: %w", err)
	}

	t.RefreshToken = nonHashedRefreshToken
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
	nonHashedRefreshToken := tokens.RefreshToken
	tokens.RefreshToken = hasher.HashString(tokens.RefreshToken)

	// save tokens in database
	err = u.repo.UpdateTokens(ctx, userID, tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to update  tokens: %w", err)
	}

	tokens.RefreshToken = nonHashedRefreshToken
	return tokens, nil
}
