package auth

import (
	"bytes"
	"client/config"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zalando/go-keyring"
)

/*
Авторизация:
 - TUI пытается взять ключи из key ring системы пользователя
  Ключи не найден:
    - делается запрос на сервер для получения новых ключей
	* Успешная авторизация
  Ключи найдены:
    - происходит проверка времени жизни найденного ключа
      Ключ просрочен:
        - делается запрос на сервер для получения новых ключей
      Ключ валиден:
		* Успешная авторизация

*/

type JWT struct {
	token        string
	refreshToken string
	expiresAt    int64
}

func GetFromLocalStorage() (*JWT, error) {
	token, err := keyring.Get(config.App.KeyRingNames.JWTAccessToken, config.App.KeyRingNames.UserName)
	if err != nil {
		return nil, err
	}

	refreshToker, err := keyring.Get(config.App.KeyRingNames.JWTRefreshToken, config.App.KeyRingNames.UserName)
	if err != nil {
		return nil, err
	}

	expiresAtString, err := keyring.Get(config.App.KeyRingNames.JWTExpiresAt, config.App.KeyRingNames.UserName)
	if err != nil {
		return nil, err
	}

	expiresAt, err := strconv.ParseInt(expiresAtString, 10, 64)
	jwtToken := &JWT{
		token:        token,
		refreshToken: refreshToker,
		expiresAt:    expiresAt,
	}
	return jwtToken, nil
}

func GetNew() (*JWT, error) {
	var claims jwt.RegisteredClaims
	var t JWT

	// do request, get link from config
	buf := bytes.Buffer{} // fixme
	tokenString := strings.TrimSpace(string(buf.Bytes()))
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, &claims)
	if err != nil {
		return nil, err
	}
	t.token = tokenString
	t.expiresAt = claims.ExpiresAt.Time.Unix()

	return &t, nil
}

func (t *JWT) Token() string {
	return t.token
}

func (t *JWT) Save() error {
	err := keyring.Set(config.App.KeyRingNames.JWTAccessToken, config.App.KeyRingNames.UserName, t.token)
	if err != nil {
		return err
	}

	err = keyring.Set(config.App.KeyRingNames.JWTRefreshToken, config.App.KeyRingNames.UserName, t.refreshToken)
	if err != nil {
		return err
	}

	err = keyring.Set(config.App.KeyRingNames.JWTExpiresAt, config.App.KeyRingNames.UserName, strconv.FormatInt(t.expiresAt, 10))
	if err != nil {
		return err
	}
	return nil
}

func (t *JWT) Expired() bool {
	return t.expiresAt-time.Now().Unix() <= 0
}
