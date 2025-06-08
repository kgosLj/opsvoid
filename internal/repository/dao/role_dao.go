package dao

import (
	"github.com/kgosLj/opsvoid/internal/model"
	"gorm.io/gorm"
)

type RoleDao interface {
	FindRoleByName(name string) (model.Role, error)
	CreateRole(role *model.Role) error
}

type GORMRoleDAO struct {
	db *gorm.DB
}

func NewRoleDao(db *gorm.DB) RoleDao {
	return &GORMRoleDAO{db: db}
}

func (dao *GORMRoleDAO) FindRoleByName(name string) (model.Role, error) {
	var role model.Role
	result := dao.db.Where("name = ?", name).First(&role)
	if result.Error == gorm.ErrRecordNotFound {
		return model.Role{}, nil
	} else if result.Error != nil {
		return model.Role{}, result.Error
	}
	return role, nil
}

func (dao *GORMRoleDAO) CreateRole(role *model.Role) error {
	result := dao.db.Create(role)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
