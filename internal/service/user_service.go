package service

import (
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository"
)

type UserService interface {
	Login(request model.LoginRequest) (model.LoginResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo
	}
}

func (svc *userService) Login(request model.LoginRequest) (model.LoginResponse, error) {
	u, err := svc.repo.FindByUsername(request)
	if err != nil {
		return model.LoginResponse{}, err
	}

}
