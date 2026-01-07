package token

import (
	"server/internal/app/config"
	domain "server/internal/app/domain/user"
	"strconv"

	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshTokenCmd struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    time.Time
	UserID       int64 `json:"user_id"`
}

func CreateNewJWT(userID int64) (string, error) {
	expirationTime := time.Now().Add(config.App.GetJWTLifetime() * time.Minute)
	userIDString := strconv.FormatInt(userID, 10)

	// Create the claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.App.GetIssuer(),
			Subject:   userIDString,
		},
	}

	// Declare the token with the signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key to get the complete encoded token string
	tokenString, err := token.SignedString([]byte(config.App.GetJWTSecret()))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return tokenString, nil
}

func VerifyJWT(jwtString string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected jwt signing algo")
		}
		return []byte(config.App.GetJWTSecret()), nil
	}

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(jwtString, claims, keyFunc,
		jwt.WithIssuer(config.App.GetIssuer()),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}
	if claims.UserID == 0 {
		return nil, domain.ErrTokenNotValid
	}
	return claims, nil
}
