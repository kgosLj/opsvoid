package dao

import (
	"github.com/kgosLj/opsvoid/internal/model"
	"gorm.io/gorm"
)

// UserDao 用户 DAO
type UserDao interface {
	FindByUsername(username string) (model.User, error)
}

// GORMUserDAO 基于 GORM 的 UserDao 实现
type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDAO{
		db: db,
	}
}

// FindByUsername 根据用户名查找用户
func (dao *GORMUserDAO) FindByUsername(username string) (model.User, error) {
	var user model.User
	err := dao.db.Preload("Role").Where("username = ?", username).First(&user).Error
	return user, err
}
