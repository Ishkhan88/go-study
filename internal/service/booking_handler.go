package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

// BookingCreate godoc
// @Summary Create booking
// @Tags bookings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param booking body model.Booking true "Booking"
// @Success 201 {object} model.Booking
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/booking [post]
func BookingCreate(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}

	var b model.Booking
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		writeErr(w, 400, "bad json")
		return
	}

	if b.ID == 0 || b.UserID == 0 || b.ConcertID == 0 {
		writeErr(w, 400, "required: id, user_id, concert_id")
		return
	}

	if b.Status == "" {
		b.Status = model.StatusPending
	}

	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	if err := repository.AddBooking(b); err != nil {
		writeErr(w, 400, err.Error())
		return
	}

	writeJSON(w, 201, b)
}

// BookingGet godoc
// @Summary Get booking by id
// @Tags bookings
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} model.Booking
// @Failure 404 {object} map[string]string
// @Router /api/booking/{id} [get]
func BookingGet(w http.ResponseWriter, r *http.Request, id int) {
	b, found := repository.GetBookingByID(id)
	if !found {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, b)
}

// BookingUpdate godoc
// @Summary Update booking by id
// @Tags bookings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param booking body model.Booking true "Booking"
// @Success 200 {object} model.Booking
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/booking/{id} [put]
func BookingUpdate(w http.ResponseWriter, r *http.Request, id int) {
	if !requireAuth(w, r) {
		return
	}

	old, found := repository.GetBookingByID(id)
	if !found {
		writeErr(w, 404, "not found")
		return
	}

	var upd model.Booking
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		writeErr(w, 400, "bad json")
		return
	}

	upd.ID = id

	if upd.UserID == 0 || upd.ConcertID == 0 {
		writeErr(w, 400, "required: user_id, concert_id")
		return
	}

	switch upd.Status {
	case model.StatusPending, model.StatusConfirmed, model.StatusRejected:
	default:
		writeErr(w, 400, "invalid status")
		return
	}

	upd.CreatedAt = old.CreatedAt
	upd.UpdatedAt = time.Now()

	updatedBooking, ok, err := repository.UpdateBooking(id, upd)
	if err != nil {
		writeErr(w, 500, "update failed")
		return
	}
	if !ok {
		writeErr(w, 404, "not found")
		return
	}

	writeJSON(w, 200, updatedBooking)
}

// BookingDelete godoc
// @Summary Delete booking by id
// @Tags bookings
// @Security BearerAuth
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/booking/{id} [delete]
func BookingDelete(w http.ResponseWriter, r *http.Request, id int) {
	if !requireAuth(w, r) {
		return
	}

	deleted, err := repository.DeleteBooking(id)
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
