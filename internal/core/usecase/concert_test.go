package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Ishkhan88/go-study/internal/apperr"
	"github.com/Ishkhan88/go-study/internal/model"
)

var (
	errConcertNextID = errors.New("concert next id failed")
	errConcertAdd    = errors.New("concert add failed")
	errConcertGet    = errors.New("concert get failed")
	errConcertUpdate = errors.New("concert update failed")
	errConcertDelete = errors.New("concert delete failed")
	errConcertList   = errors.New("concert list failed")
)

type fakeConcertRepo struct {
	nextIDVal int
	nextIDErr error

	addErr error

	byID   map[int]model.Concert
	getErr error

	updateErr      error
	updateNotFound bool

	deleteErr      error
	deleteNotFound bool

	listErr error
	list    []model.Concert
}

func (r *fakeConcertRepo) NextID(ctx context.Context) (int, error) {
	_ = ctx
	if r.nextIDErr != nil {
		return 0, r.nextIDErr
	}
	if r.nextIDVal == 0 {
		return 1, nil
	}
	return r.nextIDVal, nil
}

func (r *fakeConcertRepo) Add(ctx context.Context, c model.Concert) error {
	_ = ctx
	if r.addErr != nil {
		return r.addErr
	}
	if r.byID == nil {
		r.byID = map[int]model.Concert{}
	}
	r.byID[c.ID] = c
	return nil
}

func (r *fakeConcertRepo) GetByID(ctx context.Context, id int) (model.Concert, bool, error) {
	_ = ctx
	if r.getErr != nil {
		return model.Concert{}, false, r.getErr
	}
	if r.byID == nil {
		return model.Concert{}, false, nil
	}
	c, ok := r.byID[id]
	return c, ok, nil
}

func (r *fakeConcertRepo) Update(ctx context.Context, id int, upd model.Concert) (model.Concert, bool, error) {
	_ = ctx
	if r.updateErr != nil {
		return model.Concert{}, false, r.updateErr
	}
	if r.updateNotFound {
		return model.Concert{}, false, nil
	}
	upd.ID = id
	if r.byID == nil {
		r.byID = map[int]model.Concert{}
	}
	r.byID[id] = upd
	return upd, true, nil
}

func (r *fakeConcertRepo) Delete(ctx context.Context, id int) (bool, error) {
	_ = ctx
	if r.deleteErr != nil {
		return false, r.deleteErr
	}
	if r.deleteNotFound {
		return false, nil
	}
	if r.byID != nil {
		delete(r.byID, id)
	}
	return true, nil
}

func (r *fakeConcertRepo) List(ctx context.Context) ([]model.Concert, error) {
	_ = ctx
	if r.listErr != nil {
		return nil, r.listErr
	}
	if r.list != nil {
		return r.list, nil
	}
	// если не задан list — строим из byID
	if r.byID == nil {
		return []model.Concert{}, nil
	}
	out := make([]model.Concert, 0, len(r.byID))
	for _, v := range r.byID {
		out = append(out, v)
	}
	return out, nil
}

func TestConcertUsecase_Create(t *testing.T) {
	fixedNow := time.Date(2026, 2, 27, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		in        model.Concert
		repo      *fakeConcertRepo
		want      model.Concert
		wantErrIs error
	}{
		{
			name: "ok: generates id, sets timestamps, sets tickets_left",
			in: model.Concert{
				Title: "Show", Location: "Berlin", OrganizerEmail: "org@mail.com",
				TicketsTotal: 100, TicketsLeft: 0,
			},
			repo: &fakeConcertRepo{nextIDVal: 10},
			want: model.Concert{
				ID:             10,
				Title:          "Show",
				Location:       "Berlin",
				OrganizerEmail: "org@mail.com",
				TicketsTotal:   100,
				TicketsLeft:    100,
				CreatedAt:      fixedNow,
				UpdatedAt:      fixedNow,
			},
		},
		{
			name:      "bad input: missing required fields",
			in:        model.Concert{Title: "", Location: "X", OrganizerEmail: "a@b.com"},
			repo:      &fakeConcertRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "repo error: NextID",
			in:        model.Concert{Title: "T", Location: "L", OrganizerEmail: "a@b.com"},
			repo:      &fakeConcertRepo{nextIDErr: errConcertNextID},
			wantErrIs: errConcertNextID,
		},
		{
			name:      "repo error: Add",
			in:        model.Concert{Title: "T", Location: "L", OrganizerEmail: "a@b.com"},
			repo:      &fakeConcertRepo{nextIDVal: 1, addErr: errConcertAdd},
			wantErrIs: errConcertAdd,
		},
		{
			name: "ok: keeps provided id and createdAt if set",
			in: model.Concert{
				ID: 77, Title: "T", Location: "L", OrganizerEmail: "a@b.com",
				CreatedAt: fixedNow.Add(-time.Hour),
			},
			repo: &fakeConcertRepo{},
			want: model.Concert{
				ID:             77,
				Title:          "T",
				Location:       "L",
				OrganizerEmail: "a@b.com",
				CreatedAt:      fixedNow.Add(-time.Hour),
				UpdatedAt:      fixedNow,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewConcertUsecase(tt.repo, fakeClock{now: fixedNow})

			got, err := uc.Create(context.Background(), tt.in)
			if tt.wantErrIs != nil {
				if err == nil || !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("expected error %v, got %v", tt.wantErrIs, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// сравнение (Date/TicketPrice и др. поля сравнятся автоматически)
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("want %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestConcertUsecase_Get(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		repo      *fakeConcertRepo
		want      model.Concert
		wantErrIs error
	}{
		{
			name: "ok",
			id:   5,
			repo: &fakeConcertRepo{byID: map[int]model.Concert{5: {ID: 5, Title: "T", Location: "L", OrganizerEmail: "a@b.com"}}},
			want: model.Concert{ID: 5, Title: "T", Location: "L", OrganizerEmail: "a@b.com"},
		},
		{
			name:      "bad input",
			id:        0,
			repo:      &fakeConcertRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "not found",
			id:        123,
			repo:      &fakeConcertRepo{byID: map[int]model.Concert{}},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name:      "repo error",
			id:        5,
			repo:      &fakeConcertRepo{getErr: errConcertGet},
			wantErrIs: errConcertGet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewConcertUsecase(tt.repo, fakeClock{now: time.Time{}})

			got, err := uc.Get(context.Background(), tt.id)
			if tt.wantErrIs != nil {
				if err == nil || !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("expected error %v, got %v", tt.wantErrIs, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("want %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestConcertUsecase_Update(t *testing.T) {
	fixedNow := time.Date(2026, 2, 27, 13, 0, 0, 0, time.UTC)
	oldCreated := fixedNow.Add(-24 * time.Hour)

	tests := []struct {
		name      string
		id        int
		upd       model.Concert
		repo      *fakeConcertRepo
		want      model.Concert
		wantErrIs error
	}{
		{
			name: "ok: preserves CreatedAt and sets UpdatedAt",
			id:   7,
			upd:  model.Concert{Title: "New", Location: "L", OrganizerEmail: "a@b.com"},
			repo: &fakeConcertRepo{
				byID: map[int]model.Concert{7: {ID: 7, Title: "Old", Location: "L", OrganizerEmail: "a@b.com", CreatedAt: oldCreated}},
			},
			want: model.Concert{ID: 7, Title: "New", Location: "L", OrganizerEmail: "a@b.com", CreatedAt: oldCreated, UpdatedAt: fixedNow},
		},
		{
			name:      "bad input: id<=0",
			id:        0,
			upd:       model.Concert{Title: "T", Location: "L", OrganizerEmail: "a@b.com"},
			repo:      &fakeConcertRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "not found on GetByID",
			id:        10,
			upd:       model.Concert{Title: "T", Location: "L", OrganizerEmail: "a@b.com"},
			repo:      &fakeConcertRepo{byID: map[int]model.Concert{}},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name: "bad input: required fields empty",
			id:   10,
			upd:  model.Concert{Title: "", Location: "L", OrganizerEmail: "a@b.com"},
			repo: &fakeConcertRepo{
				byID: map[int]model.Concert{10: {ID: 10, Title: "Old", Location: "L", OrganizerEmail: "a@b.com", CreatedAt: oldCreated}},
			},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name: "repo error: Update",
			id:   10,
			upd:  model.Concert{Title: "T", Location: "L", OrganizerEmail: "a@b.com"},
			repo: &fakeConcertRepo{
				byID:      map[int]model.Concert{10: {ID: 10, Title: "Old", Location: "L", OrganizerEmail: "a@b.com", CreatedAt: oldCreated}},
				updateErr: errConcertUpdate,
			},
			wantErrIs: errConcertUpdate,
		},
		{
			name: "not found on Update",
			id:   10,
			upd:  model.Concert{Title: "T", Location: "L", OrganizerEmail: "a@b.com"},
			repo: &fakeConcertRepo{
				byID:           map[int]model.Concert{10: {ID: 10, Title: "Old", Location: "L", OrganizerEmail: "a@b.com", CreatedAt: oldCreated}},
				updateNotFound: true,
			},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name:      "repo error: GetByID",
			id:        10,
			upd:       model.Concert{Title: "T", Location: "L", OrganizerEmail: "a@b.com"},
			repo:      &fakeConcertRepo{getErr: errConcertGet},
			wantErrIs: errConcertGet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewConcertUsecase(tt.repo, fakeClock{now: fixedNow})

			got, err := uc.Update(context.Background(), tt.id, tt.upd)
			if tt.wantErrIs != nil {
				if err == nil || !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("expected error %v, got %v", tt.wantErrIs, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("want %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestConcertUsecase_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		repo      *fakeConcertRepo
		wantErrIs error
	}{
		{name: "ok", id: 1, repo: &fakeConcertRepo{}},
		{name: "bad input", id: 0, repo: &fakeConcertRepo{}, wantErrIs: apperr.ErrBadInput},
		{name: "not found", id: 99, repo: &fakeConcertRepo{deleteNotFound: true}, wantErrIs: apperr.ErrNotFound},
		{name: "repo error", id: 99, repo: &fakeConcertRepo{deleteErr: errConcertDelete}, wantErrIs: errConcertDelete},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewConcertUsecase(tt.repo, fakeClock{now: time.Time{}})

			err := uc.Delete(context.Background(), tt.id)
			if tt.wantErrIs != nil {
				if err == nil || !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("expected error %v, got %v", tt.wantErrIs, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestConcertUsecase_List(t *testing.T) {
	tests := []struct {
		name      string
		repo      *fakeConcertRepo
		wantLen   int
		wantErrIs error
	}{
		{name: "ok empty", repo: &fakeConcertRepo{list: []model.Concert{}}, wantLen: 0},
		{name: "ok non-empty", repo: &fakeConcertRepo{list: []model.Concert{{ID: 1}, {ID: 2}}}, wantLen: 2},
		{name: "repo error", repo: &fakeConcertRepo{listErr: errConcertList}, wantErrIs: errConcertList},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewConcertUsecase(tt.repo, fakeClock{now: time.Time{}})

			got, err := uc.List(context.Background())
			if tt.wantErrIs != nil {
				if err == nil || !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("expected error %v, got %v", tt.wantErrIs, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("want len %d, got %d", tt.wantLen, len(got))
			}
		})
	}
}
