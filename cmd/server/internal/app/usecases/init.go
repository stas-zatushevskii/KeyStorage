package usecases

import (
	"context"
	domain "server/internal/app/domain/user"
	"server/internal/app/usecases/user"
	"server/internal/pkg/token"
)


type UserRepository interface {
	CreateNewUser(ctx context.Context, user *domain.User) (int64, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetTokens(ctx context.Context, userId int64) (*token.Tokens, error)
	AddTokens(ctx context.Context, userId int64, token *token.Tokens) error
	UpdateTokens(ctx context.Context, userId int64, token *token.Tokens) error
}

type UseCases struct {
	UserRepository UserRepository
}

func New(userRepository UserRepository) UseCases {
	userUseCases := user.New(userRepository)
	return UseCases{
		UserRepository: userRepository,
	}
}
