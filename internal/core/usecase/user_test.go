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
	errUserNextID = errors.New("user next id failed")
	errUserAdd    = errors.New("user add failed")
	errUserGet    = errors.New("user get failed")
	errUserUpdate = errors.New("user update failed")
	errUserDelete = errors.New("user delete failed")
	errUserList   = errors.New("user list failed")
)

type fakeUserRepo struct {
	nextIDVal int
	nextIDErr error

	addErr error

	byID   map[int]model.User
	getErr error

	updateErr      error
	updateNotFound bool

	deleteErr      error
	deleteNotFound bool

	listErr error
	list    []model.User
}

func (r *fakeUserRepo) NextID(ctx context.Context) (int, error) {
	_ = ctx
	if r.nextIDErr != nil {
		return 0, r.nextIDErr
	}
	if r.nextIDVal == 0 {
		return 1, nil
	}
	return r.nextIDVal, nil
}

func (r *fakeUserRepo) Add(ctx context.Context, u model.User) error {
	_ = ctx
	if r.addErr != nil {
		return r.addErr
	}
	if r.byID == nil {
		r.byID = map[int]model.User{}
	}
	r.byID[u.ID] = u
	return nil
}

func (r *fakeUserRepo) GetByID(ctx context.Context, id int) (model.User, bool, error) {
	_ = ctx
	if r.getErr != nil {
		return model.User{}, false, r.getErr
	}
	if r.byID == nil {
		return model.User{}, false, nil
	}
	u, ok := r.byID[id]
	return u, ok, nil
}

func (r *fakeUserRepo) Update(ctx context.Context, id int, upd model.User) (model.User, bool, error) {
	_ = ctx
	if r.updateErr != nil {
		return model.User{}, false, r.updateErr
	}
	if r.updateNotFound {
		return model.User{}, false, nil
	}
	upd.ID = id
	if r.byID == nil {
		r.byID = map[int]model.User{}
	}
	r.byID[id] = upd
	return upd, true, nil
}

func (r *fakeUserRepo) Delete(ctx context.Context, id int) (bool, error) {
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

func (r *fakeUserRepo) List(ctx context.Context) ([]model.User, error) {
	_ = ctx
	if r.listErr != nil {
		return nil, r.listErr
	}
	if r.list != nil {
		return r.list, nil
	}
	if r.byID == nil {
		return []model.User{}, nil
	}
	out := make([]model.User, 0, len(r.byID))
	for _, v := range r.byID {
		out = append(out, v)
	}
	return out, nil
}

func TestUserUsecase_Create(t *testing.T) {
	fixedNow := time.Date(2026, 2, 27, 14, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		in        model.User
		repo      *fakeUserRepo
		want      model.User
		wantErrIs error
	}{
		{
			name: "ok: generates id, sets timestamps",
			in:   model.User{FirstName: "Ivan", Email: "i@mail.com"},
			repo: &fakeUserRepo{nextIDVal: 10},
			want: model.User{
				ID:        10,
				FirstName: "Ivan",
				Email:     "i@mail.com",
				CreatedAt: fixedNow,
				UpdatedAt: fixedNow,
			},
		},
		{
			name:      "bad input: missing firstname",
			in:        model.User{FirstName: "", Email: "i@mail.com"},
			repo:      &fakeUserRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "bad input: missing email",
			in:        model.User{FirstName: "Ivan", Email: ""},
			repo:      &fakeUserRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "repo error: NextID",
			in:        model.User{FirstName: "Ivan", Email: "i@mail.com"},
			repo:      &fakeUserRepo{nextIDErr: errUserNextID},
			wantErrIs: errUserNextID,
		},
		{
			name:      "repo error: Add",
			in:        model.User{FirstName: "Ivan", Email: "i@mail.com"},
			repo:      &fakeUserRepo{nextIDVal: 1, addErr: errUserAdd},
			wantErrIs: errUserAdd,
		},
		{
			name: "ok: keeps provided id and createdAt if set",
			in: model.User{
				ID: 77, FirstName: "Ivan", Email: "i@mail.com",
				CreatedAt: fixedNow.Add(-time.Hour),
			},
			repo: &fakeUserRepo{},
			want: model.User{
				ID:        77,
				FirstName: "Ivan",
				Email:     "i@mail.com",
				CreatedAt: fixedNow.Add(-time.Hour),
				UpdatedAt: fixedNow,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUsecase(tt.repo, fakeClock{now: fixedNow})

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
		})
	}
}

func TestUserUsecase_Get(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		repo      *fakeUserRepo
		want      model.User
		wantErrIs error
	}{
		{
			name: "ok",
			id:   5,
			repo: &fakeUserRepo{byID: map[int]model.User{5: {ID: 5, FirstName: "A", Email: "a@b.com"}}},
			want: model.User{ID: 5, FirstName: "A", Email: "a@b.com"},
		},
		{
			name:      "bad input",
			id:        0,
			repo:      &fakeUserRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "not found",
			id:        123,
			repo:      &fakeUserRepo{byID: map[int]model.User{}},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name:      "repo error",
			id:        5,
			repo:      &fakeUserRepo{getErr: errUserGet},
			wantErrIs: errUserGet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUsecase(tt.repo, fakeClock{now: time.Time{}})

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

func TestUserUsecase_Update(t *testing.T) {
	fixedNow := time.Date(2026, 2, 27, 15, 0, 0, 0, time.UTC)
	oldCreated := fixedNow.Add(-24 * time.Hour)

	tests := []struct {
		name      string
		id        int
		upd       model.User
		repo      *fakeUserRepo
		want      model.User
		wantErrIs error
	}{
		{
			name: "ok: preserves CreatedAt and sets UpdatedAt",
			id:   7,
			upd:  model.User{FirstName: "New", Email: "n@b.com"},
			repo: &fakeUserRepo{
				byID: map[int]model.User{7: {ID: 7, FirstName: "Old", Email: "o@b.com", CreatedAt: oldCreated}},
			},
			want: model.User{ID: 7, FirstName: "New", Email: "n@b.com", CreatedAt: oldCreated, UpdatedAt: fixedNow},
		},
		{
			name:      "bad input: id<=0",
			id:        0,
			upd:       model.User{FirstName: "A", Email: "a@b.com"},
			repo:      &fakeUserRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "not found on GetByID",
			id:        10,
			upd:       model.User{FirstName: "A", Email: "a@b.com"},
			repo:      &fakeUserRepo{byID: map[int]model.User{}},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name: "bad input: required fields empty",
			id:   10,
			upd:  model.User{FirstName: "", Email: "a@b.com"},
			repo: &fakeUserRepo{
				byID: map[int]model.User{10: {ID: 10, FirstName: "Old", Email: "o@b.com", CreatedAt: oldCreated}},
			},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name: "repo error: Update",
			id:   10,
			upd:  model.User{FirstName: "A", Email: "a@b.com"},
			repo: &fakeUserRepo{
				byID:      map[int]model.User{10: {ID: 10, FirstName: "Old", Email: "o@b.com", CreatedAt: oldCreated}},
				updateErr: errUserUpdate,
			},
			wantErrIs: errUserUpdate,
		},
		{
			name: "not found on Update",
			id:   10,
			upd:  model.User{FirstName: "A", Email: "a@b.com"},
			repo: &fakeUserRepo{
				byID:           map[int]model.User{10: {ID: 10, FirstName: "Old", Email: "o@b.com", CreatedAt: oldCreated}},
				updateNotFound: true,
			},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name:      "repo error: GetByID",
			id:        10,
			upd:       model.User{FirstName: "A", Email: "a@b.com"},
			repo:      &fakeUserRepo{getErr: errUserGet},
			wantErrIs: errUserGet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUsecase(tt.repo, fakeClock{now: fixedNow})

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

func TestUserUsecase_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		repo      *fakeUserRepo
		wantErrIs error
	}{
		{name: "ok", id: 1, repo: &fakeUserRepo{}},
		{name: "bad input", id: 0, repo: &fakeUserRepo{}, wantErrIs: apperr.ErrBadInput},
		{name: "not found", id: 99, repo: &fakeUserRepo{deleteNotFound: true}, wantErrIs: apperr.ErrNotFound},
		{name: "repo error", id: 99, repo: &fakeUserRepo{deleteErr: errUserDelete}, wantErrIs: errUserDelete},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUsecase(tt.repo, fakeClock{now: time.Time{}})

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

func TestUserUsecase_List(t *testing.T) {
	tests := []struct {
		name      string
		repo      *fakeUserRepo
		wantLen   int
		wantErrIs error
	}{
		{name: "ok empty", repo: &fakeUserRepo{list: []model.User{}}, wantLen: 0},
		{name: "ok non-empty", repo: &fakeUserRepo{list: []model.User{{ID: 1}, {ID: 2}}}, wantLen: 2},
		{name: "repo error", repo: &fakeUserRepo{listErr: errUserList}, wantErrIs: errUserList},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUsecase(tt.repo, fakeClock{now: time.Time{}})

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
