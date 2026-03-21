package service

import (
	"github.com/Ishkhan88/go-study/auth/internal/model"
	"github.com/Ishkhan88/go-study/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Register(name, email, password string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email, password string) (*model.User, error) {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) LoginWithToken(email, password string) (string, error) {

	user, err := s.Login(email, password)
	if err != nil {
		return "", err
	}

	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
