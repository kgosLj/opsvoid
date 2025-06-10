package service

import (
	"errors"
	"github.com/kgosLj/opsvoid/internal/integration/startup"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository"
	"github.com/kgosLj/opsvoid/pkg/utils"
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

	// 创建角色权限
	// 现在判断策略是否存在，如果存在就不添加了
	Sub := role.Name
	if utils.IsRbacPolicyExists(Sub, role.Obj, role.Act) {
		zap.L().Info("角色权限已存在", zap.String("role", role.Name))
		return nil
	} else {
		ok, _ := startup.E.AddPolicy(Sub, role.Obj, role.Act)
		if !ok {
			zap.L().Error("创建角色权限失败", zap.String("role", role.Name))
			return err
		}
		zap.L().Info("创建角色权限成功", zap.String("role", role.Name))
	}

	return nil
}
