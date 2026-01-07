package keyring

import (
	"client/internal/config"
	domain "client/internal/domain/token"

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

func GetFromLocalStorage() (*domain.Token, error) {
	token, err := keyring.Get(config.App.KeyRingNames.JWTAccessToken, config.App.KeyRingNames.UserName)
	if err != nil {
		return nil, err
	}

	refreshToker, err := keyring.Get(config.App.KeyRingNames.JWTRefreshToken, config.App.KeyRingNames.UserName)
	if err != nil {
		return nil, err
	}

	_, err = keyring.Get(config.App.KeyRingNames.JWTExpiresAt, config.App.KeyRingNames.UserName) // fixme
	if err != nil {
		return nil, err
	}

	tokens := domain.NewToken()
	err = tokens.SetJWTToken(token)
	if err != nil {
		return nil, err
	}
	tokens.SetRefreshToken(refreshToker)

	return tokens, nil
}

//
//func Save(t *domain.Token) error {
//	err := keyring.Set(config.App.KeyRingNames.JWTAccessToken, config.App.KeyRingNames.UserName, t.tokenJWT)
//	if err != nil {
//		return err
//	}
//
//	err = keyring.Set(config.App.KeyRingNames.JWTRefreshToken, config.App.KeyRingNames.UserName, t.tokenRefresh)
//	if err != nil {
//		return err
//	}
//
//	err = keyring.Set(config.App.KeyRingNames.JWTExpiresAt, config.App.KeyRingNames.UserName, strconv.FormatInt(t.expiresAt, 10))
//	if err != nil {
//		return err
//	}
//	return nil
//}
