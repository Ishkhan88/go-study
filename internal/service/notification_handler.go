package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

// NotificationCreate godoc
// @Summary Create notification
// @Tags notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param notification body model.Notification true "Notification"
// @Success 201 {object} model.Notification
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/notification [post]
func NotificationCreate(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}

	var n model.Notification
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeErr(w, 400, "bad json")
		return
	}

	if n.ID == 0 || n.ConcertID == 0 || n.UserID == 0 {
		writeErr(w, 400, "required: id, concert_id, user_id")
		return
	}

	if n.Status == "" {
		n.Status = model.StatusSuccess
	}

	if n.SentAt.IsZero() {
		n.SentAt = time.Now()
	}

	if err := repository.AddNotification(n); err != nil {
		writeErr(w, 400, err.Error())
		return
	}

	writeJSON(w, 201, n)
}

// NotificationGet godoc
// @Summary Get notification by id
// @Tags notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} model.Notification
// @Failure 404 {object} map[string]string
// @Router /api/notification/{id} [get]
func NotificationGet(w http.ResponseWriter, r *http.Request, id int) {
	n, found := repository.GetNotificationByID(id)
	if !found {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, n)
}

// NotificationUpdate godoc
// @Summary Update notification by id
// @Tags notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Param notification body model.Notification true "Notification"
// @Success 200 {object} model.Notification
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/notification/{id} [put]
func NotificationUpdate(w http.ResponseWriter, r *http.Request, id int) {
	if !requireAuth(w, r) {
		return
	}

	old, found := repository.GetNotificationByID(id)
	if !found {
		writeErr(w, 404, "not found")
		return
	}

	var upd model.Notification
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		writeErr(w, 400, "bad json")
		return
	}

	upd.ID = id

	if upd.ConcertID == 0 || upd.UserID == 0 {
		writeErr(w, 400, "required: concert_id, user_id")
		return
	}

	switch upd.Status {
	case model.StatusSuccess, model.StatusFailed:
	default:
		writeErr(w, 400, "invalid status")
		return
	}

	if upd.SentAt.IsZero() {
		upd.SentAt = old.SentAt
	}

	updated, ok, err := repository.UpdateNotification(id, upd)
	if err != nil {
		writeErr(w, 500, "update failed")
		return
	}
	if !ok {
		writeErr(w, 404, "not found")
		return
	}

	writeJSON(w, 200, updated)
}

// NotificationDelete godoc
// @Summary Delete notification by id
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/notification/{id} [delete]
func NotificationDelete(w http.ResponseWriter, r *http.Request, id int) {
	if !requireAuth(w, r) {
		return
	}

	deleted, err := repository.DeleteNotification(id)
	if err != nil {
		writeErr(w, 500, "delete failed")
		return
	}
	if !deleted {
		writeErr(w, 404, "not found")
		return
	}

	writeJSON(w, 200, map[string]any{
		"deleted": true,
		"id":      id,
	})
}
