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
	errNotifNextID = errors.New("notif next id failed")
	errNotifAdd    = errors.New("notif add failed")
	errNotifGet    = errors.New("notif get failed")
	errNotifUpdate = errors.New("notif update failed")
	errNotifDelete = errors.New("notif delete failed")
	errNotifList   = errors.New("notif list failed")
)

type fakeNotificationRepo struct {
	nextIDVal int
	nextIDErr error

	addErr error

	byID   map[int]model.Notification
	getErr error

	updateErr      error
	updateNotFound bool

	deleteErr      error
	deleteNotFound bool

	listErr error
	list    []model.Notification
}

func (r *fakeNotificationRepo) NextID(ctx context.Context) (int, error) {
	_ = ctx
	if r.nextIDErr != nil {
		return 0, r.nextIDErr
	}
	if r.nextIDVal == 0 {
		return 1, nil
	}
	return r.nextIDVal, nil
}

func (r *fakeNotificationRepo) Add(ctx context.Context, n model.Notification) error {
	_ = ctx
	if r.addErr != nil {
		return r.addErr
	}
	if r.byID == nil {
		r.byID = map[int]model.Notification{}
	}
	r.byID[n.ID] = n
	return nil
}

func (r *fakeNotificationRepo) GetByID(ctx context.Context, id int) (model.Notification, bool, error) {
	_ = ctx
	if r.getErr != nil {
		return model.Notification{}, false, r.getErr
	}
	if r.byID == nil {
		return model.Notification{}, false, nil
	}
	n, ok := r.byID[id]
	return n, ok, nil
}

func (r *fakeNotificationRepo) Update(ctx context.Context, id int, upd model.Notification) (model.Notification, bool, error) {
	_ = ctx
	if r.updateErr != nil {
		return model.Notification{}, false, r.updateErr
	}
	if r.updateNotFound {
		return model.Notification{}, false, nil
	}
	upd.ID = id
	if r.byID == nil {
		r.byID = map[int]model.Notification{}
	}
	r.byID[id] = upd
	return upd, true, nil
}

func (r *fakeNotificationRepo) Delete(ctx context.Context, id int) (bool, error) {
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

func (r *fakeNotificationRepo) List(ctx context.Context) ([]model.Notification, error) {
	_ = ctx
	if r.listErr != nil {
		return nil, r.listErr
	}
	if r.list != nil {
		return r.list, nil
	}
	if r.byID == nil {
		return []model.Notification{}, nil
	}
	out := make([]model.Notification, 0, len(r.byID))
	for _, v := range r.byID {
		out = append(out, v)
	}
	return out, nil
}

func TestNotificationUsecase_Create(t *testing.T) {
	fixedNow := time.Date(2026, 2, 27, 16, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		in        model.Notification
		repo      *fakeNotificationRepo
		want      model.Notification
		wantErrIs error
	}{
		{
			name: "ok: generates id when ID=0",
			in:   model.Notification{ID: 0},
			repo: &fakeNotificationRepo{nextIDVal: 5},
			want: model.Notification{ID: 5},
		},
		{
			name: "ok: keeps provided id when ID>0",
			in:   model.Notification{ID: 7},
			repo: &fakeNotificationRepo{},
			want: model.Notification{ID: 7},
		},
		{
			name:      "bad input: negative id",
			in:        model.Notification{ID: -1},
			repo:      &fakeNotificationRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "repo error: NextID",
			in:        model.Notification{ID: 0},
			repo:      &fakeNotificationRepo{nextIDErr: errNotifNextID},
			wantErrIs: errNotifNextID,
		},
		{
			name:      "repo error: Add",
			in:        model.Notification{ID: 0},
			repo:      &fakeNotificationRepo{nextIDVal: 1, addErr: errNotifAdd},
			wantErrIs: errNotifAdd,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewNotificationUsecase(tt.repo, fakeClock{now: fixedNow})

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

func TestNotificationUsecase_Get(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		repo      *fakeNotificationRepo
		want      model.Notification
		wantErrIs error
	}{
		{
			name: "ok",
			id:   1,
			repo: &fakeNotificationRepo{byID: map[int]model.Notification{1: {ID: 1}}},
			want: model.Notification{ID: 1},
		},
		{
			name:      "bad input",
			id:        0,
			repo:      &fakeNotificationRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "not found",
			id:        77,
			repo:      &fakeNotificationRepo{byID: map[int]model.Notification{}},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name:      "repo error",
			id:        1,
			repo:      &fakeNotificationRepo{getErr: errNotifGet},
			wantErrIs: errNotifGet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewNotificationUsecase(tt.repo, fakeClock{now: time.Time{}})

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

func TestNotificationUsecase_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		upd       model.Notification
		repo      *fakeNotificationRepo
		want      model.Notification
		wantErrIs error
	}{
		{
			name: "ok: sets ID to key",
			id:   7,
			upd:  model.Notification{ID: 999}, // должен перезаписаться на 7
			repo: &fakeNotificationRepo{
				byID: map[int]model.Notification{7: {ID: 7}},
			},
			want: model.Notification{ID: 7},
		},
		{
			name:      "bad input: id<=0",
			id:        0,
			upd:       model.Notification{ID: 1},
			repo:      &fakeNotificationRepo{},
			wantErrIs: apperr.ErrBadInput,
		},
		{
			name:      "not found on GetByID",
			id:        10,
			upd:       model.Notification{ID: 1},
			repo:      &fakeNotificationRepo{byID: map[int]model.Notification{}},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name: "repo error: Update",
			id:   10,
			upd:  model.Notification{ID: 1},
			repo: &fakeNotificationRepo{
				byID:      map[int]model.Notification{10: {ID: 10}},
				updateErr: errNotifUpdate,
			},
			wantErrIs: errNotifUpdate,
		},
		{
			name: "not found on Update",
			id:   10,
			upd:  model.Notification{ID: 1},
			repo: &fakeNotificationRepo{
				byID:           map[int]model.Notification{10: {ID: 10}},
				updateNotFound: true,
			},
			wantErrIs: apperr.ErrNotFound,
		},
		{
			name:      "repo error: GetByID",
			id:        10,
			upd:       model.Notification{ID: 1},
			repo:      &fakeNotificationRepo{getErr: errNotifGet},
			wantErrIs: errNotifGet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewNotificationUsecase(tt.repo, fakeClock{now: time.Time{}})

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

func TestNotificationUsecase_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		repo      *fakeNotificationRepo
		wantErrIs error
	}{
		{name: "ok", id: 1, repo: &fakeNotificationRepo{}},
		{name: "bad input", id: 0, repo: &fakeNotificationRepo{}, wantErrIs: apperr.ErrBadInput},
		{name: "not found", id: 99, repo: &fakeNotificationRepo{deleteNotFound: true}, wantErrIs: apperr.ErrNotFound},
		{name: "repo error", id: 99, repo: &fakeNotificationRepo{deleteErr: errNotifDelete}, wantErrIs: errNotifDelete},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewNotificationUsecase(tt.repo, fakeClock{now: time.Time{}})

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

func TestNotificationUsecase_List(t *testing.T) {
	tests := []struct {
		name      string
		repo      *fakeNotificationRepo
		wantLen   int
		wantErrIs error
	}{
		{name: "ok empty", repo: &fakeNotificationRepo{list: []model.Notification{}}, wantLen: 0},
		{name: "ok non-empty", repo: &fakeNotificationRepo{list: []model.Notification{{ID: 1}, {ID: 2}}}, wantLen: 2},
		{name: "repo error", repo: &fakeNotificationRepo{listErr: errNotifList}, wantErrIs: errNotifList},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			uc := NewNotificationUsecase(tt.repo, fakeClock{now: time.Time{}})

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
