package service

import (
	"errors"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository"
	"go.uber.org/zap"
)

var (
	ExistRole = errors.New("角色已存在")
)

type RoleService interface {
	CreateRole(role *model.RoleCreateRequest) error
}

type UserRoleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &UserRoleService{repo: repo}
}

func (s *UserRoleService) CreateRole(role *model.RoleCreateRequest) error {
	dbrole, err := s.repo.FindRoleByName(role.Name)
	if err == nil && dbrole.Name != "" {
		zap.L().Error("角色已存在", zap.String("role", dbrole.Name))
		return ExistRole
	}
	zap.L().Info("创建角色", zap.String("role", role.Name))
	err = s.repo.CreateRole(role)
	if err != nil {
		zap.L().Error("创建角色失败", zap.Error(err), zap.String("role", role.Name))
		return err
	}
	zap.L().Info("创建角色成功", zap.String("role", role.Name))
	return nil
}
