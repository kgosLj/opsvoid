package dao

import (
	"github.com/kgosLj/opsvoid/internal/model"
	"gorm.io/gorm"
)

type UserDao interface {
	FindByUsername(request model.LoginRequest) (model.User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDAO{db: db}
}
