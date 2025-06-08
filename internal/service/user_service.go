package service

import (
	"errors"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository"
	"github.com/kgosLj/opsvoid/internal/web/middleware/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrUserPassword = errors.New("密码错误")
)

type UserService interface {
	Login(request model.LoginRequest) (model.LoginResponse, error)
}

// userService 实现 UserService 接口
type userService struct {
	repo repository.UserRepository
}

// NewUserService 构造函数
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// Login 登录功能逻辑
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
		zap.L().Error("密码比较失败", zap.Error(err), zap.String("username", username), zap.String("stored_password_hash", u.Password), zap.String("provided_password", password))
		return model.LoginResponse{}, ErrUserPassword
	}

	// 如果密码正确就生成 token
	token := jwt.GenerateToken(u)

	return model.LoginResponse{
		Username: u.Username,
		Token:    token,
	}, nil
}
