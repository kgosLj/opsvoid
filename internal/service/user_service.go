package service

import (
	"errors"
	"fmt"
	"github.com/kgosLj/opsvoid/internal/integration/startup"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository"
	"github.com/kgosLj/opsvoid/internal/web/middleware/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound  = errors.New("用户不存在")
	ErrUserPassword  = errors.New("密码错误")
	ExistUser        = errors.New("用户已存在")
	ErrHashPassword  = errors.New("密码加密失败")
	NotFoundRole     = errors.New("角色不存在")
	NotFoundPassword = errors.New("密码不能为空")
)

type UserService interface {
	Login(request model.LoginRequest) (model.LoginResponse, error)               // 登录功能
	GetUserInfo(username string) (model.GetUserInfo, error)                      // 获得自身的用户信息
	CreateUser(user *model.CreateUserRequest) (*model.CreateUserResponse, error) // 创建新用户
	BindRole(req *model.BindRoleRequest) (*model.BindRoleResponse, error)        // 用户绑定角色
}

// userService 实现 UserService 接口
type userService struct {
	repo repository.UserRepository
}

// NewUserService 构造函数
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// Login 登录功能逻辑
func (svc *userService) Login(request model.LoginRequest) (model.LoginResponse, error) {
	username := request.Username
	password := request.Password
	u, err := svc.repo.FindByUsername(username)
	if errors.Is(err, ErrUserNotFound) {
		return model.LoginResponse{}, ErrUserNotFound
	}
	if err != nil {
		return model.LoginResponse{}, err
	}

	// 判断密码是否有正确
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		zap.L().Error("密码比较失败", zap.Error(err), zap.String("username", username), zap.String("stored_password_hash", u.Password), zap.String("provided_password", password))
		return model.LoginResponse{}, ErrUserPassword
	}

	// 如果密码正确就生成 token
	token := jwt.GenerateToken(u)

	return model.LoginResponse{
		Username: u.Username,
		Token:    token,
	}, nil
}

// GetUserInfo 获取用户信息
func (svc *userService) GetUserInfo(username string) (model.GetUserInfo, error) {
	user, err := svc.repo.FindByUsername(username)
	if errors.Is(err, ErrUserNotFound) {
		return model.GetUserInfo{}, ErrUserNotFound
	} else if err != nil {
		return model.GetUserInfo{}, err
	}
	return model.GetUserInfo{
		Username: user.Username,
		Role:     user.Role[0].Name,
	}, nil
}

// CreateUser 创建用户
func (svc *userService) CreateUser(req *model.CreateUserRequest) (*model.CreateUserResponse, error) {
	// 校验用户名是否存在
	_, err := svc.repo.FindByUsername(req.Username)
	if err == nil {
		return &model.CreateUserResponse{}, ExistUser
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return &model.CreateUserResponse{}, fmt.Errorf("查询用户失败: %v", err) // 数据库报错
	}

	if req.Password == "" {
		return &model.CreateUserResponse{}, NotFoundPassword
	}
	// 哈希密码
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &model.CreateUserResponse{}, ErrHashPassword
	}

	// 查询角色
	var roles []*model.Role
	for _, roleName := range req.Role {
		role, err := svc.repo.FindRoleByName(roleName)
		if err != nil {
			return &model.CreateUserResponse{}, NotFoundRole
		}
		roles = append(roles, &role)
	}

	// 创建用户对象
	dbUser := &model.User{
		Username: req.Username,
		Password: string(hashPassword),
		Role:     roles,
	}

	if err := svc.repo.CreateUser(dbUser); err != nil {
		zap.L().Error("创建用户失败", zap.Error(err), zap.String("username", req.Username))
		return &model.CreateUserResponse{}, err
	}

	// fix: 同步更新 casbin 的用户-角色关系策略 (g)
	for _, roleName := range req.Role {
		_, err = startup.E.AddGroupingPolicy(req.Username, roleName)
		if err != nil {
			zap.L().Error("casbin 添加用户-角色策略失败", zap.Error(err), zap.String("username", req.Username), zap.String("role", roleName))
			return nil, err
		}
		zap.L().Info("casbin 添加用户-角色策略成功", zap.String("username", req.Username), zap.String("role", roleName))
	}

	return &model.CreateUserResponse{
		Username: dbUser.Username,
		Role:     req.Role,
	}, nil
}

// BindRole 绑定用户权限 (同时需要注意的是，要同时使用 enforce 来绑定上用户名和 casbin 的关系)
func (svc *userService) BindRole(req *model.BindRoleRequest) (*model.BindRoleResponse, error) {
	// 首先先判断是否存在 用户 和 角色
	// 1. 判断用户是否存在
	user, err := svc.repo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 2. 检查角色是否存在
	var roles []*model.Role
	for _, roleName := range req.Roles {
		role, err := svc.repo.FindRoleByName(roleName)
		if err != nil {
			if errors.Is(err, repository.ErrRoleNotFound) {
				return nil, NotFoundRole
			}
			return nil, fmt.Errorf("查询角色失败: %v", err)
		}
		newRole := role
		roles = append(roles, &newRole)
	}

	//3. 绑定角色（假设通过 GORM 多对多关系关联更新）
	if err := svc.repo.UpdateUserRoles(&user, roles); err != nil {
		zap.L().Error("更新用户角色失败", zap.Error(err), zap.String("username", req.Username))
		return nil, err
	}

	// fix: 同步更新 casbin 的用户-角色关系策略 (g)
	for _, roleName := range req.Roles {
		_, err = startup.E.AddGroupingPolicy(user.Username, roleName)
		if err != nil {
			zap.L().Error("casbin 添加用户-角色策略失败", zap.Error(err), zap.String("username", user.Username), zap.String("role", roleName))
			return nil, fmt.Errorf("casbin 添加用户-角色策略失败: %v", err)
		}
	}

	// 返回响应
	return &model.BindRoleResponse{
		Username: user.Username,
		Roles:    req.Roles,
	}, nil

}
