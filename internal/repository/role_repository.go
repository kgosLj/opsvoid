package repository

import (
	"errors"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository/dao"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindRoleByName(name string) (model.Role, error)
	CreateRole(role *model.RoleCreateRequest) error
}

type CachedRoleRepository struct {
	dao dao.RoleDao
}

func NewRoleRepository(dao dao.RoleDao) RoleRepository {
	return &CachedRoleRepository{dao: dao}
}

func (r *CachedRoleRepository) FindRoleByName(name string) (model.Role, error) {
	role, err := r.dao.FindRoleByName(name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Role{}, ErrRoleNotFound
	} else if err != nil {
		return model.Role{}, err
	}
	return role, nil
}
func (r *CachedRoleRepository) CreateRole(role *model.RoleCreateRequest) error {
	dbrole := &model.Role{
		Desc: role.Desc,
		Name: role.Name,
	}
	return r.dao.CreateRole(dbrole)
}
