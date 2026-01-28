package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

// UserCreate godoc
// @Summary Create user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body model.User true "User"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/user [post]
func UserCreate(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}

	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeErr(w, 400, "bad json")
		return
	}

	if u.ID == 0 || u.FirstName == "" || u.Email == "" {
		writeErr(w, 400, "required: id, first_name, email")
		return
	}

	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now

	if err := repository.AddUser(u); err != nil {
		writeErr(w, 400, err.Error())
		return
	}

	writeJSON(w, 201, u)
}

// UserGet godoc
// @Summary Get user by id
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 404 {object} map[string]string
// @Router /api/user/{id} [get]
func UserGet(w http.ResponseWriter, r *http.Request, id int) {
	u, found := repository.GetUserByID(id)
	if !found {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, u)
}

// UserUpdate godoc
// @Summary Update user by id
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body model.User true "User"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/user/{id} [put]
func UserUpdate(w http.ResponseWriter, r *http.Request, id int) {
	if !requireAuth(w, r) {
		return
	}

	old, found := repository.GetUserByID(id)
	if !found {
		writeErr(w, 404, "not found")
		return
	}

	var upd model.User
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		writeErr(w, 400, "bad json")
		return
	}

	upd.ID = id

	if upd.FirstName == "" || upd.Email == "" {
		writeErr(w, 400, "required: first_name, email")
		return
	}

	upd.CreatedAt = old.CreatedAt
	upd.UpdatedAt = time.Now()

	updatedUser, ok, err := repository.UpdateUser(id, upd)
	if err != nil {
		writeErr(w, 500, "update failed")
		return
	}
	if !ok {
		writeErr(w, 404, "not found")
		return
	}

	writeJSON(w, 200, updatedUser)
}

// UserDelete godoc
// @Summary Delete user by id
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/user/{id} [delete]
func UserDelete(w http.ResponseWriter, r *http.Request, id int) {
	if !requireAuth(w, r) {
		return
	}

	deleted, err := repository.DeleteUser(id)
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
