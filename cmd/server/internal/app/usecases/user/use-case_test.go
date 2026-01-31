package user

import (
	"context"
	"errors"
	"server/internal/app/config"
	domain "server/internal/app/domain/user"
	hasher "server/internal/pkg/hash/argon2"
	"server/internal/pkg/token"
	"strings"
	"testing"
	"time"
)

func init() {
	config.InitTestConfig()
}

type repoFake struct {
	createNewUser func(ctx context.Context, user *domain.User) (int64, error)
	getByUsername func(ctx context.Context, username string) (*domain.User, error)
	getTokens     func(ctx context.Context, userId int64) (*token.Tokens, error)
	addTokens     func(ctx context.Context, userId int64, t *token.Tokens) error
	updateTokens  func(ctx context.Context, userId int64, t *token.Tokens) error
}

func (r *repoFake) CreateNewUser(ctx context.Context, user *domain.User) (int64, error) {
	if r.createNewUser != nil {
		return r.createNewUser(ctx, user)
	}
	return 0, nil
}
func (r *repoFake) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if r.getByUsername != nil {
		return r.getByUsername(ctx, username)
	}
	return nil, nil
}
func (r *repoFake) GetTokens(ctx context.Context, userId int64) (*token.Tokens, error) {
	if r.getTokens != nil {
		return r.getTokens(ctx, userId)
	}
	return nil, nil
}
func (r *repoFake) AddTokens(ctx context.Context, userId int64, t *token.Tokens) error {
	if r.addTokens != nil {
		return r.addTokens(ctx, userId, t)
	}
	return nil
}
func (r *repoFake) UpdateTokens(ctx context.Context, userId int64, t *token.Tokens) error {
	if r.updateTokens != nil {
		return r.updateTokens(ctx, userId, t)
	}
	return nil
}

func TestUser_RefreshJWTToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid jwt -> returns verify error (not mapped)", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.RefreshJWTToken(ctx, "bad-jwt", "refresh")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("repo.GetTokens error -> returned as-is", func(t *testing.T) {
		t.Parallel()

		jwt, err := token.CreateNewJWT(1)
		if err != nil {
			t.Fatalf("CreateNewJWT error: %v", err)
		}

		dbErr := errors.New("db down")
		uc := New(&repoFake{
			getTokens: func(ctx context.Context, userId int64) (*token.Tokens, error) {
				return nil, dbErr
			},
		})

		_, err = uc.RefreshJWTToken(ctx, jwt, "refresh")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected wrapped/returned dbErr, got: %v", err)
		}
	})

	t.Run("tokens revoked -> ErrTokenRevoked", func(t *testing.T) {
		t.Parallel()

		jwt, err := token.CreateNewJWT(10)
		if err != nil {
			t.Fatalf("CreateNewJWT error: %v", err)
		}

		uc := New(&repoFake{
			getTokens: func(ctx context.Context, userId int64) (*token.Tokens, error) {
				return &token.Tokens{
					UserId:            userId,
					Revoked:           true,
					RefreshTokenExpAt: time.Now().Add(10 * time.Minute),
					RefreshToken:      "hash-doesnt-matter",
				}, nil
			},
		})

		_, err = uc.RefreshJWTToken(ctx, jwt, "any")
		if !errors.Is(err, domain.ErrTokenRevoked) {
			t.Fatalf("expected ErrTokenRevoked, got: %v", err)
		}
	})

	t.Run("refresh expired -> ErrRefreshTokenExpired", func(t *testing.T) {
		t.Parallel()

		jwt, err := token.CreateNewJWT(11)
		if err != nil {
			t.Fatalf("CreateNewJWT error: %v", err)
		}

		uc := New(&repoFake{
			getTokens: func(ctx context.Context, userId int64) (*token.Tokens, error) {
				return &token.Tokens{
					UserId:            userId,
					Revoked:           false,
					RefreshTokenExpAt: time.Now().Add(-1 * time.Minute),
					RefreshToken:      "hash-doesnt-matter",
				}, nil
			},
		})

		_, err = uc.RefreshJWTToken(ctx, jwt, "any")
		if !errors.Is(err, domain.ErrRefreshTokenExpired) {
			t.Fatalf("expected ErrRefreshTokenExpired, got: %v", err)
		}
	})

	t.Run("invalid refresh token -> ErrInvalidRefreshToken", func(t *testing.T) {
		t.Parallel()

		jwt, err := token.CreateNewJWT(12)
		if err != nil {
			t.Fatalf("CreateNewJWT error: %v", err)
		}

		hashed, err := hasher.HashString("right-refresh")
		if err != nil {
			t.Fatalf("HashString error: %v", err)
		}

		uc := New(&repoFake{
			getTokens: func(ctx context.Context, userId int64) (*token.Tokens, error) {
				return &token.Tokens{
					UserId:            userId,
					Revoked:           false,
					RefreshTokenExpAt: time.Now().Add(10 * time.Minute),
					RefreshToken:      hashed,
				}, nil
			},
		})

		_, err = uc.RefreshJWTToken(ctx, jwt, "wrong-refresh")
		if !errors.Is(err, domain.ErrInvalidRefreshToken) {
			t.Fatalf("expected ErrInvalidRefreshToken, got: %v", err)
		}
	})

	t.Run("verify refresh returns error -> wrapped", func(t *testing.T) {
		t.Parallel()

		jwt, err := token.CreateNewJWT(13)
		if err != nil {
			t.Fatalf("CreateNewJWT error: %v", err)
		}

		badHash := "not-a-valid-hash-format"

		uc := New(&repoFake{
			getTokens: func(ctx context.Context, userId int64) (*token.Tokens, error) {
				return &token.Tokens{
					UserId:            userId,
					Revoked:           false,
					RefreshTokenExpAt: time.Now().Add(10 * time.Minute),
					RefreshToken:      badHash,
				}, nil
			},
		})

		_, err = uc.RefreshJWTToken(ctx, jwt, "any-refresh")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to verify refresh token") {
			t.Fatalf("expected wrapped verify error, got: %v", err)
		}
	})

	t.Run("ok -> UpdateTokens called, returns new tokens", func(t *testing.T) {
		t.Parallel()

		userId := int64(1001)

		jwt, err := token.CreateNewJWT(userId)
		if err != nil {
			t.Fatalf("CreateNewJWT error: %v", err)
		}

		hashed, err := hasher.HashString("refresh-plain")
		if err != nil {
			t.Fatalf("HashString error: %v", err)
		}

		updateCalled := false

		uc := New(&repoFake{
			getTokens: func(ctx context.Context, uid int64) (*token.Tokens, error) {
				if uid != userId {
					t.Fatalf("expected userId=%d, got %d", userId, uid)
				}
				return &token.Tokens{
					UserId:            uid,
					Revoked:           false,
					RefreshTokenExpAt: time.Now().Add(10 * time.Minute),
					RefreshToken:      hashed,
				}, nil
			},
			updateTokens: func(ctx context.Context, uid int64, tks *token.Tokens) error {
				updateCalled = true
				if uid != userId {
					t.Fatalf("expected uid=%d, got %d", userId, uid)
				}
				if tks.RefreshToken == "" {
					t.Fatalf("expected non-empty refresh hash")
				}
				if tks.RefreshToken == "refresh-plain" {
					t.Fatalf("expected hashed refresh in repo, got plain")
				}
				if tks.JWTToken == "" {
					t.Fatalf("expected jwt set")
				}
				return nil
			},
		})

		newTokens, err := uc.RefreshJWTToken(ctx, jwt, "refresh-plain")
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if !updateCalled {
			t.Fatalf("expected UpdateTokens to be called")
		}

		if newTokens.RefreshToken == "" {
			t.Fatalf("expected plain refresh token in response")
		}
		if newTokens.JWTToken == "" {
			t.Fatalf("expected jwt in response")
		}
	})
}
