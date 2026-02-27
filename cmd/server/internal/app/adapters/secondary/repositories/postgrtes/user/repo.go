package user

import (
	"context"
	"database/sql"
	"errors"
	"server/internal/app/domain/user"
	"server/internal/pkg/token"

	"github.com/jackc/pgconn"
)

func (u *Repository) GetById(ctx context.Context, id int64) (*user.User, error) {
	query := `
		SELECT id, username, password_hash
		FROM users
		WHERE id = $1`

	newUser := user.NewUser()
	if err := u.db.QueryRowContext(ctx, query, id).Scan(&newUser.ID, &newUser.Username, &newUser.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return newUser, nil
}

func (u *Repository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
		SELECT id, username, password_hash
		FROM users
		WHERE username = $1`

	newUser := user.NewUser()
	if err := u.db.QueryRowContext(ctx, query, username).Scan(&newUser.ID, &newUser.Username, &newUser.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return newUser, nil
}

func (u *Repository) CreateNewUser(ctx context.Context, newUser *user.User) (int64, error) {
	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id`

	var id int64
	err := u.db.QueryRowContext(ctx, query, newUser.Username, newUser.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return 0, user.ErrUsernameAlreadyExists
			}
		}
		return 0, err
	}
	return id, nil
}

func (u *Repository) AddTokens(ctx context.Context, userId int64, token *token.Tokens) error {
	query := `
		INSERT INTO user_tokens (user_id, refresh_token, refresh_token_expires_at, revoked_at)
		VALUES ($1, $2, $3, $4)`
	_, err := u.db.ExecContext(ctx, query, userId, token.RefreshToken, token.RefreshTokenExpAt, nil)
	if err != nil {
		return err
	}
	return nil
}

func (u *Repository) UpdateTokens(ctx context.Context, userId int64, token *token.Tokens) error {
	query := `
		UPDATE user_tokens
		SET 
			refresh_token = $2,
			refresh_token_expires_at = $3,
			revoked_at = $4
		WHERE user_id = $1`
	_, err := u.db.ExecContext(ctx, query, userId, token.RefreshToken, token.RefreshTokenExpAt, nil)
	if err != nil {
		return err
	}
	return nil
}

func (u *Repository) GetTokens(ctx context.Context, userId int64) (*token.Tokens, error) {
	query := `
		SELECT refresh_token, refresh_token_expires_at, revoked_at
		FROM user_tokens
		WHERE user_id = $1`

	var revokedAt sql.NullInt64
	newToken := token.NewTokens(userId)

	err := u.db.QueryRowContext(ctx, query, userId).Scan(&newToken.RefreshToken, &newToken.RefreshTokenExpAt, &revokedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrRefreshTokenNotFound
		}
		return nil, err
	}
	if revokedAt.Valid {
		newToken.Revoked = true
	}
	return newToken, err
}
