package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

// ConcertCreate godoc
// @Summary Create concert
// @Tags concerts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param concert body model.Concert true "Concert"
// @Success 201 {object} model.Concert
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/concert [post]
func ConcertCreate(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}

	var c model.Concert
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeErr(w, 400, "bad json")
		return
	}

	// Подстрой под свои обязательные поля (ниже пример логики как у User)
	if c.ID == 0 || c.Title == "" || c.Location == "" || c.OrganizerEmail == "" {
		writeErr(w, 400, "required: id, title, location, organizer_email")
		return
	}

	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	if err := repository.AddConcert(c); err != nil {
		writeErr(w, 400, err.Error())
		return
	}

	writeJSON(w, 201, c)
}

// ConcertGet godoc
// @Summary Get concert by id
// @Tags concerts
// @Produce json
// @Param id path int true "Concert ID"
// @Success 200 {object} model.Concert
// @Failure 404 {object} map[string]string
// @Router /api/concert/{id} [get]
func ConcertGet(w http.ResponseWriter, r *http.Request, id int) {
	c, found := repository.GetConcertByID(id)
	if !found {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, c)
}

// ConcertUpdate godoc
// @Summary Update concert by id
// @Tags concerts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Concert ID"
// @Param concert body model.Concert true "Concert"
// @Success 200 {object} model.Concert
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/concert/{id} [put]
func ConcertUpdate(w http.ResponseWriter, r *http.Request, id int) {
	if !requireAuth(w, r) {
		return
	}

	old, found := repository.GetConcertByID(id)
	if !found {
		writeErr(w, 404, "not found")
		return
	}

	var upd model.Concert
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		writeErr(w, 400, "bad json")
		return
	}

	upd.ID = id

	if upd.Title == "" || upd.Location == "" || upd.OrganizerEmail == "" {
		writeErr(w, 400, "required: title, location, organizer_email")
		return
	}

	upd.CreatedAt = old.CreatedAt
	upd.UpdatedAt = time.Now()

	updatedConcert, ok, err := repository.UpdateConcert(id, upd)
	if err != nil {
		writeErr(w, 500, "update failed")
		return
	}
	if !ok {
		writeErr(w, 404, "not found")
		return
	}

	writeJSON(w, 200, updatedConcert)
}

// ConcertDelete godoc
// @Summary Delete concert by id
// @Tags concerts
// @Security BearerAuth
// @Produce json
// @Param id path int true "Concert ID"
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/concert/{id} [delete]
func ConcertDelete(w http.ResponseWriter, r *http.Request, id int) {
	if !requireAuth(w, r) {
		return
	}

	deleted, err := repository.DeleteConcert(id)
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
