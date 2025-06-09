package dao

import (
	"github.com/kgosLj/opsvoid/internal/model"
	"gorm.io/gorm"
)

// UserDao 用户 DAO
type UserDao interface {
	FindByUsername(username string) (model.User, error)
	CreateUser(user *model.User) error
	FindRoleByName(name string) (model.Role, error)
	UpdateUserRoles(user *model.User, roles []*model.Role) error
}

// GORMUserDAO 基于 GORM 的 UserDao 实现
type GORMUserDAO struct {
	db *gorm.DB
}

// NewUserDao 创建一个新的 UserDao 实例
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

// CreateUser 创建用户
func (dao *GORMUserDAO) CreateUser(user *model.User) error {
	return dao.db.Create(user).Error
}

// FindRoleByName 根据角色名称查找角色
func (dao *GORMUserDAO) FindRoleByName(name string) (model.Role, error) {
	var role model.Role
	err := dao.db.Where("name =?", name).First(&role).Error
	return role, err
}

func (dao *GORMUserDAO) UpdateUserRoles(user *model.User, roles []*model.Role) error {
	return dao.db.Model(user).Association("Role").Replace(roles)
}
