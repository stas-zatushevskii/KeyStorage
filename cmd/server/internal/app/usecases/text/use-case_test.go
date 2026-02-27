package text

import (
	"context"
	"errors"
	"testing"

	domain "server/internal/app/domain/text_obj"
)

type repoFake struct {
	getByUserID func(ctx context.Context, userId int64) ([]*domain.Text, error)
	getByID     func(ctx context.Context, textId int64) (*domain.Text, error)
	create      func(ctx context.Context, text *domain.Text) (int64, error)
	update      func(ctx context.Context, text *domain.Text) error
}

func (r *repoFake) GetByUserID(ctx context.Context, userId int64) ([]*domain.Text, error) {
	if r.getByUserID != nil {
		return r.getByUserID(ctx, userId)
	}
	return nil, nil
}
func (r *repoFake) GetByID(ctx context.Context, textId int64) (*domain.Text, error) {
	if r.getByID != nil {
		return r.getByID(ctx, textId)
	}
	return nil, nil
}
func (r *repoFake) Create(ctx context.Context, text *domain.Text) (int64, error) {
	if r.create != nil {
		return r.create(ctx, text)
	}
	return 0, nil
}
func (r *repoFake) Update(ctx context.Context, text *domain.Text) error {
	if r.update != nil {
		return r.update(ctx, text)
	}
	return nil
}

func TestTextObj_GetText(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid textId -> ErrInvalidTextID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.GetText(ctx, 0)
		if !errors.Is(err, domain.ErrInvalidTextID) {
			t.Fatalf("expected ErrInvalidTextID, got: %v", err)
		}
	})

	t.Run("repo error -> ErrTextNotFound", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByID: func(ctx context.Context, textId int64) (*domain.Text, error) {
				return nil, errors.New("db error")
			},
		})

		_, err := uc.GetText(ctx, 1)
		if !errors.Is(err, domain.ErrTextNotFound) {
			t.Fatalf("expected ErrTextNotFound, got: %v", err)
		}
	})

	t.Run("ok -> returns item", func(t *testing.T) {
		t.Parallel()

		want := &domain.Text{TextId: 7, UserId: 1, Title: "t", Text: "body"}

		uc := New(&repoFake{
			getByID: func(ctx context.Context, textId int64) (*domain.Text, error) {
				return want, nil
			},
		})

		got, err := uc.GetText(ctx, 7)
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if got != want {
			t.Fatalf("expected same pointer")
		}
	})
}

func TestTextObj_GetTextList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid userId -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.GetTextList(ctx, 0)
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("repo error -> ErrTextNotFound", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByUserID: func(ctx context.Context, userId int64) ([]*domain.Text, error) {
				return nil, errors.New("db error")
			},
		})

		_, err := uc.GetTextList(ctx, 1)
		if !errors.Is(err, domain.ErrTextNotFound) {
			t.Fatalf("expected ErrTextNotFound, got: %v", err)
		}
	})

	t.Run("empty list -> ErrEmptyTextsList", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			getByUserID: func(ctx context.Context, userId int64) ([]*domain.Text, error) {
				return []*domain.Text{}, nil
			},
		})

		_, err := uc.GetTextList(ctx, 1)
		if !errors.Is(err, domain.ErrEmptyTextsList) {
			t.Fatalf("expected ErrEmptyTextsList, got: %v", err)
		}
	})

	t.Run("ok -> returns list", func(t *testing.T) {
		t.Parallel()

		want := []*domain.Text{
			{TextId: 1, UserId: 1, Title: "a", Text: "x"},
			{TextId: 2, UserId: 1, Title: "b", Text: "y"},
		}

		uc := New(&repoFake{
			getByUserID: func(ctx context.Context, userId int64) ([]*domain.Text, error) {
				return want, nil
			},
		})

		got, err := uc.GetTextList(ctx, 1)
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if len(got) != len(want) {
			t.Fatalf("expected len=%d, got len=%d", len(want), len(got))
		}
		if got[0] != want[0] || got[1] != want[1] {
			t.Fatalf("expected same pointers in slice")
		}
	})
}

func TestTextObj_CreateNewTextObj(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("nil text -> ErrFailedCreateText", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewTextObj(ctx, nil)
		if !errors.Is(err, domain.ErrFailedCreateText) {
			t.Fatalf("expected ErrFailedCreateText, got: %v", err)
		}
	})

	t.Run("invalid userId -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewTextObj(ctx, &domain.Text{UserId: 0, Title: "t", Text: "x"})
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("empty title -> ErrEmptyTitle", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewTextObj(ctx, &domain.Text{UserId: 1, Title: "", Text: "x"})
		if !errors.Is(err, domain.ErrEmptyTitle) {
			t.Fatalf("expected ErrEmptyTitle, got: %v", err)
		}
	})

	t.Run("empty text -> ErrEmptyText", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		_, err := uc.CreateNewTextObj(ctx, &domain.Text{UserId: 1, Title: "t", Text: ""})
		if !errors.Is(err, domain.ErrEmptyText) {
			t.Fatalf("expected ErrEmptyText, got: %v", err)
		}
	})

	t.Run("repo error -> ErrFailedCreateText", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			create: func(ctx context.Context, text *domain.Text) (int64, error) {
				return 0, errors.New("insert failed")
			},
		})

		_, err := uc.CreateNewTextObj(ctx, &domain.Text{UserId: 1, Title: "t", Text: "x"})
		if !errors.Is(err, domain.ErrFailedCreateText) {
			t.Fatalf("expected ErrFailedCreateText, got: %v", err)
		}
	})

	t.Run("ok -> returns id", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			create: func(ctx context.Context, text *domain.Text) (int64, error) {
				if text.UserId != 1 {
					t.Fatalf("expected UserId=1, got %d", text.UserId)
				}
				if text.Title != "hello" {
					t.Fatalf("expected Title=hello, got %q", text.Title)
				}
				if text.Text != "world" {
					t.Fatalf("expected Text=world, got %q", text.Text)
				}
				return 55, nil
			},
		})

		id, err := uc.CreateNewTextObj(ctx, &domain.Text{UserId: 1, Title: "hello", Text: "world"})
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if id != 55 {
			t.Fatalf("expected id=55, got %d", id)
		}
	})
}

func TestTextObj_UpdateText(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("nil text -> ErrFailedUpdateText", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateText(ctx, nil)
		if !errors.Is(err, domain.ErrFailedUpdateText) {
			t.Fatalf("expected ErrFailedUpdateText, got: %v", err)
		}
	})

	t.Run("invalid userId -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateText(ctx, &domain.Text{UserId: 0, Title: "t", Text: "x"})
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("empty title -> ErrEmptyTitle", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateText(ctx, &domain.Text{UserId: 1, Title: "", Text: "x"})
		if !errors.Is(err, domain.ErrEmptyTitle) {
			t.Fatalf("expected ErrEmptyTitle, got: %v", err)
		}
	})

	t.Run("empty text -> ErrEmptyText", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{})
		err := uc.UpdateText(ctx, &domain.Text{UserId: 1, Title: "t", Text: ""})
		if !errors.Is(err, domain.ErrEmptyText) {
			t.Fatalf("expected ErrEmptyText, got: %v", err)
		}
	})

	t.Run("repo error -> ErrFailedUpdateText", func(t *testing.T) {
		t.Parallel()

		uc := New(&repoFake{
			update: func(ctx context.Context, text *domain.Text) error {
				return errors.New("update failed")
			},
		})

		err := uc.UpdateText(ctx, &domain.Text{UserId: 1, Title: "t", Text: "x"})
		if !errors.Is(err, domain.ErrFailedUpdateText) {
			t.Fatalf("expected ErrFailedUpdateText, got: %v", err)
		}
	})

	t.Run("ok -> nil", func(t *testing.T) {
		t.Parallel()

		called := false

		uc := New(&repoFake{
			update: func(ctx context.Context, text *domain.Text) error {
				called = true
				if text.UserId != 1 {
					t.Fatalf("expected UserId=1, got %d", text.UserId)
				}
				if text.Title != "title" {
					t.Fatalf("expected Title=title, got %q", text.Title)
				}
				if text.Text != "body" {
					t.Fatalf("expected Text=body, got %q", text.Text)
				}
				return nil
			},
		})

		err := uc.UpdateText(ctx, &domain.Text{UserId: 1, Title: "title", Text: "body"})
		if err != nil {
			t.Fatalf("expected nil err, got: %v", err)
		}
		if !called {
			t.Fatalf("expected repo.Update to be called")
		}
	})
}
