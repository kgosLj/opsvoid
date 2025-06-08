package repository

import (
	"errors"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository/dao"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
)

type UserRepository interface {
	FindByUsername(username string) (model.User, error)
}

type CachedUserRepository struct {
	dao dao.UserDao
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &CachedUserRepository{
		dao: dao,
	}
}

func (r *CachedUserRepository) FindByUsername(username string) (model.User, error) {
	user, err := r.dao.FindByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
