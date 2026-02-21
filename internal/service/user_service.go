package service

import (
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func CreateUser(u model.User) (model.User, error) {
	// Create не должен требовать ID от клиента
	if u.FirstName == "" || u.Email == "" {
		return model.User{}, ErrBadInput
	}

	// если ID не передали — сгенерируем
	if u.ID == 0 {
		u.ID = repository.GetNextUserID()
	}

	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now

	if err := repository.AddUser(u); err != nil {
		return model.User{}, err
	}

	return u, nil
}

func GetUser(id int) (model.User, error) {
	u, ok := repository.GetUserByID(id)
	if !ok {
		return model.User{}, ErrNotFound
	}
	return u, nil
}

func UpdateUser(id int, upd model.User) (model.User, error) {
	old, ok := repository.GetUserByID(id)
	if !ok {
		return model.User{}, ErrNotFound
	}

	// обязательные поля
	if upd.FirstName == "" || upd.Email == "" {
		return model.User{}, ErrBadInput
	}

	upd.ID = id
	upd.CreatedAt = old.CreatedAt
	upd.UpdatedAt = time.Now()

	updated, ok2, err := repository.UpdateUser(id, upd)
	if err != nil {
		return model.User{}, err
	}
	if !ok2 {
		return model.User{}, ErrNotFound
	}

	return updated, nil
}

func DeleteUser(id int) error {
	deleted, err := repository.DeleteUser(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func ListUsers() []model.User {
	return repository.GetUserSafeCopy()
}
