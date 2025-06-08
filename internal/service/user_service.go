package service

import (
	"errors"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrUserPassword = errors.New("密码错误")
)

type UserService interface {
	Login(request model.LoginRequest) (model.LoginResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Login(request model.LoginRequest) (model.LoginResponse, error) {
	username := request.Username
	password := request.Password
	u, err := svc.repo.FindByUsername(username)
	if errors.Is(err, ErrUserNotFound) {
		return model.LoginResponse{}, ErrUserNotFound
	}
	if err != nil {
		return model.LoginResponse{}, err
	}

	// 判断密码是否有正确
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return model.LoginResponse{}, ErrUserPassword
	}

	// 如果密码正确就生成 token
	return model.LoginResponse{
		Username: u.Username,
		Token:
	}, nil

}
