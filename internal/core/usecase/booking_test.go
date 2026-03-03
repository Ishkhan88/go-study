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
	errNextID = errors.New("next id failed")
	errAdd    = errors.New("add failed")
	errGet    = errors.New("get failed")
	errUpdate = errors.New("update failed")
	errDelete = errors.New("delete failed")
	errList   = errors.New("list failed")
)

// ---- fakes ----

type fakeClock struct {
	now time.Time
}

func (c fakeClock) Now() time.Time { return c.now }

// fakeBookingRepo реализует port.BookingRepository (через duck-typing на методы).
type fakeBookingRepo struct {
	nextIDVal int
	nextIDErr error

	addErr error
	added  []model.Booking

	byID    map[int]model.Booking
	getErr  error
	listErr error

	updateErr error
	deleteErr error

	// Можно управлять not-found для update/delete отдельно
	updateNotFound bool
	deleteNotFound bool
}

func (r *fakeBookingRepo) NextID(ctx context.Context) (int, error) {
	_ = ctx
	if r.nextIDErr != nil {
		return 0, r.nextIDErr
	}
	if r.nextIDVal == 0 {
		return 1, nil
	}
	return r.nextIDVal, nil
}

func (r *fakeBookingRepo) Add(ctx context.Context, b model.Booking) error {
	_ = ctx
	if r.addErr != nil {
		return r.addErr
	}
	r.added = append(r.added, b)
	if r.byID == nil {
		r.byID = map[int]model.Booking{}
	}
	r.byID[b.ID] = b
	return nil
}

func (r *fakeBookingRepo) GetByID(ctx context.Context, id int) (model.Booking, bool, error) {
	_ = ctx
	if r.getErr != nil {
		return model.Booking{}, false, r.getErr
	}
	if r.byID == nil {
		return model.Booking{}, false, nil
	}
	b, ok := r.byID[id]
	return b, ok, nil
}

func (r *fakeBookingRepo) Update(ctx context.Context, id int, upd model.Booking) (model.Booking, bool, error) {
	_ = ctx
	if r.updateErr != nil {
		return model.Booking{}, false, r.updateErr
	}
	if r.updateNotFound {
		return model.Booking{}, false, nil
	}
	upd.ID = id
	if r.byID == nil {
		r.byID = map[int]model.Booking{}
	}
	r.byID[id] = upd
	return upd, true, nil
}

func (r *fakeBookingRepo) Delete(ctx context.Context, id int) (bool, error) {
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

func (r *fakeBookingRepo) List(ctx context.Context) ([]model.Booking, error) {
	_ = ctx
	if r.listErr != nil {
		return nil, r.listErr
	}
	if r.byID == nil {
		return []model.Booking{}, nil
	}
	out := make([]model.Booking, 0, len(r.byID))
	for _, v := range r.byID {
		out = append(out, v)
	}
	return out, nil
}

// ---- tests ----

func TestBookingUsecase_Create(t *testing.T) {
	fixedNow := time.Date(2026, 2, 27, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		in        model.Booking
		repo      *fakeBookingRepo
		want      model.Booking
		wantErrIs error
	}{
		{
			name: "ok: assigns id, status, timestamps",
			in:   model.Booking{UserID: 10, ConcertID: 20, Status: ""},
			repo: &fakeBookingRepo{nextIDVal: 7},
			want: model.Booking{
				ID:        7,
				UserID:    10,
				ConcertID: 20,
				Status:    model.StatusPending,
				CreatedAt: fixedNow,
				UpdatedAt: fixedNow,
			},
		},
		{
			name:      "bad input: user id <= 0",
			in:        model.Booking{UserID: 0, ConcertID: 20},
			repo:      &fakeBookingRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "bad input: concert id <= 0",
			in:        model.Booking{UserID: 10, ConcertID: 0},
			repo:      &fakeBookingRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "repo error: NextID",
			in:        model.Booking{UserID: 10, ConcertID: 20},
			repo:      &fakeBookingRepo{nextIDErr: errNextID},
			wantErrIs: errNextID,
		},
		{
			name:      "repo error: Add",
			in:        model.Booking{UserID: 10, ConcertID: 20},
			repo:      &fakeBookingRepo{nextIDVal: 3, addErr: errAdd},
			wantErrIs: errAdd,
		},
		{
			name: "ok: keeps provided status",
			in:   model.Booking{UserID: 10, ConcertID: 20, Status: model.StatusConfirmed},
			repo: &fakeBookingRepo{nextIDVal: 9},
			want: model.Booking{
				ID:        9,
				UserID:    10,
				ConcertID: 20,
				Status:    model.StatusConfirmed,
				CreatedAt: fixedNow,
				UpdatedAt: fixedNow,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewBookingUsecase(tt.repo, fakeClock{now: fixedNow})

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

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("want %+v, got %+v", tt.want, got)
			}

			if len(tt.repo.added) != 1 {
				t.Fatalf("expected 1 added booking, got %d", len(tt.repo.added))
			}
		})
	}
}

func TestBookingUsecase_Get(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		repo      *fakeBookingRepo
		want      model.Booking
		wantErrIs error
	}{
		{
			name: "ok",
			id:   5,
			repo: &fakeBookingRepo{
				byID: map[int]model.Booking{5: {ID: 5, UserID: 1, ConcertID: 2, Status: model.StatusPending}},
			},
			want: model.Booking{ID: 5, UserID: 1, ConcertID: 2, Status: model.StatusPending},
		},
		{
			name:      "bad input",
			id:        0,
			repo:      &fakeBookingRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "not found",
			id:        123,
			repo:      &fakeBookingRepo{byID: map[int]model.Booking{}},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name:      "repo error",
			id:        5,
			repo:      &fakeBookingRepo{getErr: errGet},
			wantErrIs: errGet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewBookingUsecase(tt.repo, fakeClock{now: time.Time{}})

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

func TestBookingUsecase_Update(t *testing.T) {
	fixedNow := time.Date(2026, 2, 27, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        int
		upd       model.Booking
		repo      *fakeBookingRepo
		want      model.Booking
		wantErrIs error
	}{
		{
			name: "ok: sets UpdatedAt",
			id:   7,
			upd:  model.Booking{UserID: 2, ConcertID: 3, Status: model.StatusConfirmed},
			repo: &fakeBookingRepo{},
			want: model.Booking{ID: 7, UserID: 2, ConcertID: 3, Status: model.StatusConfirmed, UpdatedAt: fixedNow},
		},
		{
			name:      "bad input",
			id:        0,
			upd:       model.Booking{UserID: 2, ConcertID: 3},
			repo:      &fakeBookingRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "not found",
			id:        10,
			upd:       model.Booking{UserID: 2, ConcertID: 3},
			repo:      &fakeBookingRepo{updateNotFound: true},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name:      "repo error",
			id:        10,
			upd:       model.Booking{UserID: 2, ConcertID: 3},
			repo:      &fakeBookingRepo{updateErr: errUpdate},
			wantErrIs: errUpdate,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewBookingUsecase(tt.repo, fakeClock{now: fixedNow})

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

			if got.ID != tt.want.ID ||
				got.UserID != tt.want.UserID ||
				got.ConcertID != tt.want.ConcertID ||
				got.Status != tt.want.Status ||
				!got.UpdatedAt.Equal(tt.want.UpdatedAt) {
				t.Fatalf("want %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestBookingUsecase_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		repo      *fakeBookingRepo
		wantErrIs error
	}{
		{
			name: "ok",
			id:   1,
			repo: &fakeBookingRepo{},
		},
		{
			name:      "bad input",
			id:        0,
			repo:      &fakeBookingRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "not found",
			id:        99,
			repo:      &fakeBookingRepo{deleteNotFound: true},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name:      "repo error",
			id:        99,
			repo:      &fakeBookingRepo{deleteErr: errDelete},
			wantErrIs: errDelete,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewBookingUsecase(tt.repo, fakeClock{now: time.Time{}})

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

func TestBookingUsecase_List(t *testing.T) {
	tests := []struct {
		name      string
		repo      *fakeBookingRepo
		wantLen   int
		wantErrIs error
	}{
		{
			name:    "ok: empty",
			repo:    &fakeBookingRepo{byID: map[int]model.Booking{}},
			wantLen: 0,
		},
		{
			name: "ok: non-empty",
			repo: &fakeBookingRepo{
				byID: map[int]model.Booking{
					1: {ID: 1, UserID: 1, ConcertID: 1, Status: model.StatusPending},
					2: {ID: 2, UserID: 2, ConcertID: 2, Status: model.StatusConfirmed},
				},
			},
			wantLen: 2,
		},
		{
			name:      "repo error",
			repo:      &fakeBookingRepo{listErr: errList},
			wantErrIs: errList,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewBookingUsecase(tt.repo, fakeClock{now: time.Time{}})

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
