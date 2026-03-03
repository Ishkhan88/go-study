package service

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/core/port/clock"
	"github.com/Ishkhan88/go-study/internal/core/usecase"
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func userUC() usecase.UserUsecase {
	return usecase.NewUserUsecase(
		repository.UserRepoAdapter{},
		clock.SystemClock{},
	)
}

func CreateUser(u model.User) (model.User, error) {
	return userUC().Create(context.Background(), u)
}

func GetUser(id int) (model.User, error) {
	return userUC().Get(context.Background(), id)
}

func UpdateUser(id int, upd model.User) (model.User, error) {
	return userUC().Update(context.Background(), id, upd)
}

func DeleteUser(id int) error {
	return userUC().Delete(context.Background(), id)
}

func ListUsers() []model.User {
	list, err := userUC().List(context.Background())
	if err != nil {
		return []model.User{}
	}
	return list
}
