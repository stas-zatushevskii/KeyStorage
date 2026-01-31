package bank_card_obj

import (
	"context"
	"errors"
	"testing"

	domain "server/internal/app/domain/bank_card_obj"
)

type repoFake struct {
	getByUserID func(ctx context.Context, userId int64) ([]*domain.BankCard, error)
	getByID     func(ctx context.Context, cardId int64) (*domain.BankCard, error)
	create      func(ctx context.Context, card *domain.BankCard) (int64, error)
	update      func(ctx context.Context, card *domain.BankCard) error
}

func (r *repoFake) GetByUserID(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
	if r.getByUserID != nil {
		return r.getByUserID(ctx, userId)
	}
	return nil, nil
}

func (r *repoFake) GetByID(ctx context.Context, cardId int64) (*domain.BankCard, error) {
	if r.getByID != nil {
		return r.getByID(ctx, cardId)
	}
	return nil, nil
}

func (r *repoFake) Create(ctx context.Context, card *domain.BankCard) (int64, error) {
	if r.create != nil {
		return r.create(ctx, card)
	}
	return 0, nil
}

func (r *repoFake) Update(ctx context.Context, card *domain.BankCard) error {
	if r.update != nil {
		return r.update(ctx, card)
	}
	return nil
}

func TestBankCardObj_GetBankCard(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid cardId -> ErrInvalidCardID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.GetBankCard(ctx, 0)
		if !errors.Is(err, domain.ErrInvalidCardID) {
			t.Fatalf("expected ErrInvalidCardID, got: %v", err)
		}
	})

	t.Run("repo error -> ErrBankCardNotFound", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByID: func(ctx context.Context, cardId int64) (*domain.BankCard, error) {
				return nil, errors.New("db error")
			},
		})

		_, err := uc.GetBankCard(ctx, 10)
		if !errors.Is(err, domain.ErrBankCardNotFound) {
			t.Fatalf("expected ErrBankCardNotFound, got: %v", err)
		}
	})

	t.Run("ok -> returns card", func(t *testing.T) {
		t.Parallel()

		want := &domain.BankCard{CardId: 7, UserId: 1, Bank: "maib", Pid: "PID123"}

		uc := New(&repoFake{
			getByID: func(ctx context.Context, cardId int64) (*domain.BankCard, error) {
				return want, nil
			},
		})

		got, err := uc.GetBankCard(ctx, 7)
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if got != want {
			t.Fatalf("expected same pointer, got different")
		}
	})
}

func TestBankCardObj_GetBankCardList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid userId -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.GetBankCardList(ctx, -1)
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("repo error -> ErrBankCardNotFound", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByUserID: func(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
				return nil, errors.New("select failed")
			},
		})

		_, err := uc.GetBankCardList(ctx, 1)
		if !errors.Is(err, domain.ErrBankCardNotFound) {
			t.Fatalf("expected ErrBankCardNotFound, got: %v", err)
		}
	})

	t.Run("empty list -> ErrBankCardNotFound", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByUserID: func(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
				return []*domain.BankCard{}, nil
			},
		})

		_, err := uc.GetBankCardList(ctx, 1)
		if !errors.Is(err, domain.ErrBankCardNotFound) {
			t.Fatalf("expected ErrBankCardNotFound, got: %v", err)
		}
	})

	t.Run("ok -> returns list", func(t *testing.T) {
		t.Parallel()

		want := []*domain.BankCard{
			{CardId: 1, UserId: 1, Bank: "maib", Pid: "PID1"},
			{CardId: 2, UserId: 1, Bank: "victoriabank", Pid: "PID2"},
		}

		uc := New(&repoFake{
			getByUserID: func(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
				return want, nil
			},
		})

		got, err := uc.GetBankCardList(ctx, 1)
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

func TestBankCardObj_CreateNewBankCardObj(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("nil card -> ErrFaildeCreateBankCardObject", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewBankCardObj(ctx, nil)
		if !errors.Is(err, domain.ErrFaildeCreateBankCardObject) {
			t.Fatalf("expected ErrFaildeCreateBankCardObject, got: %v", err)
		}
	})

	t.Run("invalid userId -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewBankCardObj(ctx, &domain.BankCard{UserId: 0, Bank: "maib", Pid: "PID"})
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("empty bank -> ErrEmptyBankName", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewBankCardObj(ctx, &domain.BankCard{UserId: 1, Bank: "", Pid: "PID"})
		if !errors.Is(err, domain.ErrEmptyBankName) {
			t.Fatalf("expected ErrEmptyBankName, got: %v", err)
		}
	})

	t.Run("empty pid -> ErrEmptyPID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewBankCardObj(ctx, &domain.BankCard{UserId: 1, Bank: "maib", Pid: ""})
		if !errors.Is(err, domain.ErrEmptyPID) {
			t.Fatalf("expected ErrEmptyPID, got: %v", err)
		}
	})

	t.Run("repo error -> ErrFaildeCreateBankCardObject", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			create: func(ctx context.Context, card *domain.BankCard) (int64, error) {
				return 0, errors.New("insert failed")
			},
		})

		_, err := uc.CreateNewBankCardObj(ctx, &domain.BankCard{UserId: 1, Bank: "maib", Pid: "PID"})
		if !errors.Is(err, domain.ErrFaildeCreateBankCardObject) {
			t.Fatalf("expected ErrFaildeCreateBankCardObject, got: %v", err)
		}
	})

	t.Run("ok -> returns id", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			create: func(ctx context.Context, card *domain.BankCard) (int64, error) {
				if card.UserId != 1 {
					t.Fatalf("expected UserId=1, got %d", card.UserId)
				}
				if card.Bank != "maib" {
					t.Fatalf("expected Bank=maib, got %q", card.Bank)
				}
				if card.Pid != "PID123" {
					t.Fatalf("expected Pid=PID123, got %q", card.Pid)
				}
				return 100, nil
			},
		})

		id, err := uc.CreateNewBankCardObj(ctx, &domain.BankCard{UserId: 1, Bank: "maib", Pid: "PID123"})
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if id != 100 {
			t.Fatalf("expected id=100, got %d", id)
		}
	})
}

func TestBankCardObj_UpdateBankCard(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("nil card -> ErrFailedUpdateBankCard", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateBankCard(ctx, nil)
		if !errors.Is(err, domain.ErrFailedUpdateBankCard) {
			t.Fatalf("expected ErrFailedUpdateBankCard, got: %v", err)
		}
	})

	t.Run("invalid userId -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateBankCard(ctx, &domain.BankCard{UserId: 0, Bank: "maib", Pid: "PID"})
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("empty bank -> ErrEmptyBankName", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateBankCard(ctx, &domain.BankCard{UserId: 1, Bank: "", Pid: "PID"})
		if !errors.Is(err, domain.ErrEmptyBankName) {
			t.Fatalf("expected ErrEmptyBankName, got: %v", err)
		}
	})

	t.Run("empty pid -> ErrEmptyPID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateBankCard(ctx, &domain.BankCard{UserId: 1, Bank: "maib", Pid: ""})
		if !errors.Is(err, domain.ErrEmptyPID) {
			t.Fatalf("expected ErrEmptyPID, got: %v", err)
		}
	})

	t.Run("repo error -> ErrFailedUpdateBankCard", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			update: func(ctx context.Context, card *domain.BankCard) error {
				return errors.New("update failed")
			},
		})

		err := uc.UpdateBankCard(ctx, &domain.BankCard{UserId: 1, Bank: "maib", Pid: "PID"})
		if !errors.Is(err, domain.ErrFailedUpdateBankCard) {
			t.Fatalf("expected ErrFailedUpdateBankCard, got: %v", err)
		}
	})

	t.Run("ok -> nil", func(t *testing.T) {
		t.Parallel()

		called := false

		uc := New(&repoFake{
			update: func(ctx context.Context, card *domain.BankCard) error {
				called = true
				if card.UserId != 1 {
					t.Fatalf("expected UserId=1, got %d", card.UserId)
				}
				if card.Bank != "maib" {
					t.Fatalf("expected Bank=maib, got %q", card.Bank)
				}
				if card.Pid != "PID1" {
					t.Fatalf("expected Pid=PID1, got %q", card.Pid)
				}
				return nil
			},
		})

		err := uc.UpdateBankCard(ctx, &domain.BankCard{UserId: 1, Bank: "maib", Pid: "PID1"})
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if !called {
			t.Fatalf("expected repo.Update to be called")
		}
	})
}
