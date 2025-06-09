package repository

import (
	"errors"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository/dao"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrRoleNotFound = errors.New("角色不存在")
)

// UserRepository 用户仓库接口
type UserRepository interface {
	FindByUsername(username string) (model.User, error)
	CreateUser(user *model.User) error
	FindRoleByName(name string) (model.Role, error)
	UpdateUserRoles(user *model.User, roles []*model.Role) error
}

type CachedUserRepository struct {
	dao dao.UserDao
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &CachedUserRepository{
		dao: dao,
	}
}

// FindByUsername 根据用户名查找用户
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

// CreateUser 创建用户
func (r *CachedUserRepository) CreateUser(user *model.User) error {
	return r.dao.CreateUser(user)
}

// FindRoleByName 根据角色名称查找角色
func (r *CachedUserRepository) FindRoleByName(name string) (model.Role, error) {
	role, err := r.dao.FindRoleByName(name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Role{}, ErrRoleNotFound
	}
	return role, nil
}

func (r *CachedUserRepository) UpdateUserRoles(user *model.User, roles []*model.Role) error {
	err := r.dao.UpdateUserRoles(user, roles)
	return err
}
