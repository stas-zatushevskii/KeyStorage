package token

import (
	"errors"
	"testing"
	"time"

	"server/internal/app/config"
	domain "server/internal/app/domain/user"

	"github.com/golang-jwt/jwt/v5"
)

func init() {
	config.InitTestConfig()
}

func requireJWTConfig(t *testing.T) {
	t.Helper()

	if config.App == nil {
		t.Skip("config.App is nil (config not initialized in tests)")
	}

	secret := config.App.GetJWTSecret()
	issuer := config.App.GetIssuer()

	if secret == "" || issuer == "" {
		t.Skipf("jwt config is not initialized (secret=%q issuer=%q)", secret, issuer)
	}
}

func TestCreateNewJWT_ThenVerify_OK(t *testing.T) {
	requireJWTConfig(t)

	userID := int64(123)

	j, err := CreateNewJWT(userID)
	if err != nil {
		t.Fatalf("CreateNewJWT error: %v", err)
	}
	if j == "" {
		t.Fatalf("expected non-empty jwt")
	}

	claims, err := VerifyJWT(j)
	if err != nil {
		t.Fatalf("VerifyJWT error: %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("expected UserID=%d, got %d", userID, claims.UserID)
	}

	if claims.Issuer != config.App.GetIssuer() {
		t.Fatalf("expected Issuer=%q, got %q", config.App.GetIssuer(), claims.Issuer)
	}
	if claims.Subject == "" {
		t.Fatalf("expected Subject to be set")
	}
	if claims.ExpiresAt == nil {
		t.Fatalf("expected ExpiresAt to be set")
	}
}

func TestVerifyJWT_InvalidString(t *testing.T) {
	requireJWTConfig(t)

	_, err := VerifyJWT("not-a-jwt")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestVerifyJWT_UserIDZero_ReturnsErrTokenNotValid(t *testing.T) {
	requireJWTConfig(t)

	j, err := CreateNewJWT(0)
	if err != nil {
		t.Fatalf("CreateNewJWT error: %v", err)
	}

	_, err = VerifyJWT(j)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrTokenNotValid) {
		t.Fatalf("expected ErrTokenNotValid, got: %v", err)
	}
}

func TestVerifyJWT_WrongAlgorithm(t *testing.T) {
	requireJWTConfig(t)

	secret := []byte(config.App.GetJWTSecret())

	claims := &Claims{
		UserID: 42,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.App.GetIssuer(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Subject:   "42",
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	j, err := tok.SignedString(secret)
	if err != nil {
		t.Fatalf("SignedString error: %v", err)
	}

	_, err = VerifyJWT(j)
	if err == nil {
		t.Fatalf("expected error for wrong signing method, got nil")
	}
}

func TestVerifyJWT_WrongIssuer(t *testing.T) {
	requireJWTConfig(t)

	secret := []byte(config.App.GetJWTSecret())

	claims := &Claims{
		UserID: 7,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "some-other-issuer",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Subject:   "7",
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	j, err := tok.SignedString(secret)
	if err != nil {
		t.Fatalf("SignedString error: %v", err)
	}

	_, err = VerifyJWT(j)
	if err == nil {
		t.Fatalf("expected error for wrong issuer, got nil")
	}
}
