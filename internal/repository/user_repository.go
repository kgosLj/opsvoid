package repository

import (
	"errors"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository/dao"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(request model.LoginRequest) (model.User, error)
}

type CachedUserRepository struct {
	dao dao.UserDao
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &CachedUserRepository{
		dao: dao,
	}
}

func (r *CachedUserRepository) FindByUsername(request model.LoginRequest) (model.User, error) {
	user, err := r.dao.FindByUsername(request)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, errors.New("not found")
	}
	if err != nil {
		return model.User{}, err
	}
	return user, nil

}
