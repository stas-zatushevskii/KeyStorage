package account

import (
	"context"
	"errors"
	"testing"

	domain "server/internal/app/domain/account_obj"
)

type repoFake struct {
	getByUserID func(ctx context.Context, userId int64) ([]*domain.Account, error)
	getByID     func(ctx context.Context, accountId int64) (*domain.Account, error)
	create      func(ctx context.Context, account *domain.Account) (int64, error)
	update      func(ctx context.Context, account *domain.Account) error
}

func (r *repoFake) GetByUserID(ctx context.Context, userId int64) ([]*domain.Account, error) {
	if r.getByUserID != nil {
		return r.getByUserID(ctx, userId)
	}
	return nil, nil
}

func (r *repoFake) GetByID(ctx context.Context, accountId int64) (*domain.Account, error) {
	if r.getByID != nil {
		return r.getByID(ctx, accountId)
	}
	return nil, nil
}

func (r *repoFake) Create(ctx context.Context, account *domain.Account) (int64, error) {
	if r.create != nil {
		return r.create(ctx, account)
	}
	return 0, nil
}

func (r *repoFake) Update(ctx context.Context, account *domain.Account) error {
	if r.update != nil {
		return r.update(ctx, account)
	}
	return nil
}

func TestAccountObj_GetAccountsList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid userId -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.GetAccountsList(ctx, 0)
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("repo error -> ErrAccountNotFound", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByUserID: func(ctx context.Context, userId int64) ([]*domain.Account, error) {
				return nil, errors.New("db down")
			},
		})

		_, err := uc.GetAccountsList(ctx, 1)
		if !errors.Is(err, domain.ErrAccountNotFound) {
			t.Fatalf("expected ErrAccountNotFound, got: %v", err)
		}
	})

	t.Run("empty list -> ErrEmptyAccountsList", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByUserID: func(ctx context.Context, userId int64) ([]*domain.Account, error) {
				return []*domain.Account{}, nil
			},
		})

		_, err := uc.GetAccountsList(ctx, 1)
		if !errors.Is(err, domain.ErrEmptyAccountsList) {
			t.Fatalf("expected ErrEmptyAccountsList, got: %v", err)
		}
	})

	t.Run("ok -> returns list", func(t *testing.T) {
		t.Parallel()

		want := []*domain.Account{
			{AccountId: 10, UserId: 1, ServiceName: "telegram"},
			{AccountId: 11, UserId: 1, ServiceName: "shopify"},
		}

		uc := New(&repoFake{
			getByUserID: func(ctx context.Context, userId int64) ([]*domain.Account, error) {
				return want, nil
			},
		})

		got, err := uc.GetAccountsList(ctx, 1)
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if len(got) != len(want) {
			t.Fatalf("expected len=%d, got len=%d", len(want), len(got))
		}
		if got[0] != want[0] || got[1] != want[1] {
			t.Fatalf("expected same slice items pointers")
		}
	})
}

func TestAccountObj_GetAccount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid accountId -> ErrInvalidAccountID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.GetAccount(ctx, -1)
		if !errors.Is(err, domain.ErrInvalidAccountID) {
			t.Fatalf("expected ErrInvalidAccountID, got: %v", err)
		}
	})

	t.Run("repo error -> ErrAccountNotFound", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByID: func(ctx context.Context, accountId int64) (*domain.Account, error) {
				return nil, errors.New("not found in db")
			},
		})

		_, err := uc.GetAccount(ctx, 123)
		if !errors.Is(err, domain.ErrAccountNotFound) {
			t.Fatalf("expected ErrAccountNotFound, got: %v", err)
		}
	})

	t.Run("ok -> returns account", func(t *testing.T) {
		t.Parallel()

		want := &domain.Account{AccountId: 7, UserId: 1, ServiceName: "amoCRM"}

		uc := New(&repoFake{
			getByID: func(ctx context.Context, accountId int64) (*domain.Account, error) {
				return want, nil
			},
		})

		got, err := uc.GetAccount(ctx, 7)
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if got != want {
			t.Fatalf("expected same pointer, got different")
		}
	})
}

func TestAccountObj_CreateNewAccountObj(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid userId -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewAccountObj(ctx, &domain.Account{UserId: 0, ServiceName: "x"})
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("empty serviceName -> ErrEmptyServiceName", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewAccountObj(ctx, &domain.Account{UserId: 1, ServiceName: ""})
		if !errors.Is(err, domain.ErrEmptyServiceName) {
			t.Fatalf("expected ErrEmptyServiceName, got: %v", err)
		}
	})

	t.Run("repo error -> ErrFailedCreateAccount", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			create: func(ctx context.Context, account *domain.Account) (int64, error) {
				return 0, errors.New("insert failed")
			},
		})

		_, err := uc.CreateNewAccountObj(ctx, &domain.Account{UserId: 1, ServiceName: "telegram"})
		if !errors.Is(err, domain.ErrFailedCreateAccount) {
			t.Fatalf("expected ErrFailedCreateAccount, got: %v", err)
		}
	})

	t.Run("ok -> returns id", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			create: func(ctx context.Context, account *domain.Account) (int64, error) {
				if account.UserId != 1 {
					t.Fatalf("expected UserId=1, got %d", account.UserId)
				}
				if account.ServiceName != "telegram" {
					t.Fatalf("expected ServiceName=telegram, got %q", account.ServiceName)
				}
				return 42, nil
			},
		})

		id, err := uc.CreateNewAccountObj(ctx, &domain.Account{UserId: 1, ServiceName: "telegram"})
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if id != 42 {
			t.Fatalf("expected id=42, got %d", id)
		}
	})
}

func TestAccountObj_UpdateAccount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid accountId -> ErrInvalidAccountID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateAccount(ctx, &domain.Account{AccountId: 0, ServiceName: "x"})
		if !errors.Is(err, domain.ErrInvalidAccountID) {
			t.Fatalf("expected ErrInvalidAccountID, got: %v", err)
		}
	})

	t.Run("empty serviceName -> ErrEmptyServiceName", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateAccount(ctx, &domain.Account{AccountId: 1, ServiceName: ""})
		if !errors.Is(err, domain.ErrEmptyServiceName) {
			t.Fatalf("expected ErrEmptyServiceName, got: %v", err)
		}
	})

	t.Run("repo error -> ErrFailedUpdateAccount", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			update: func(ctx context.Context, account *domain.Account) error {
				return errors.New("update failed")
			},
		})

		err := uc.UpdateAccount(ctx, &domain.Account{AccountId: 1, ServiceName: "telegram"})
		if !errors.Is(err, domain.ErrFailedUpdateAccount) {
			t.Fatalf("expected ErrFailedUpdateAccount, got: %v", err)
		}
	})

	t.Run("ok -> nil", func(t *testing.T) {
		t.Parallel()

		called := false

		uc := New(&repoFake{
			update: func(ctx context.Context, account *domain.Account) error {
				called = true
				if account.AccountId != 99 {
					t.Fatalf("expected AccountId=99, got %d", account.AccountId)
				}
				if account.ServiceName != "amoCRM" {
					t.Fatalf("expected ServiceName=amoCRM, got %q", account.ServiceName)
				}
				return nil
			},
		})

		err := uc.UpdateAccount(ctx, &domain.Account{AccountId: 99, ServiceName: "amoCRM"})
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if !called {
			t.Fatalf("expected repo.Update to be called")
		}
	})
}
