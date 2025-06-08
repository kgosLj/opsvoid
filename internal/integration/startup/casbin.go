package startup

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Casbin 初始化 Casbin
func InitEnforce(db *gorm.DB) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		zap.L().Fatal("初始化 Casbin 适配器失败", zap.Error(err))
	}
	e, err := casbin.NewEnforcer("./config/casbin.conf", adapter)
	if err != nil {
		zap.L().Fatal("从 Casbin 模型初始化失败", zap.Error(err))
	}

	// 从数据库加载权限策略
	err = e.LoadPolicy()
	if err != nil {
		zap.L().Error("从数据库加载 Casbin 策略失败", zap.Error(err))
	}

	// --- 添加一些默认策略 ---
	// exp: 让 "admin" 角色的用户可以操作所有 /api/v1/ 下的任何资源
	if ok, _ := e.AddPolicy("admin", "/api/v1/*", "*"); !ok {
		zap.L().Info("策略已存在: admin, /api/v1/*, *")
	}

	// exp: 让 "guest" 角色只有查看自己信息的权限
	if ok, _ := e.AddPolicy("guest", "/api/v1/user/info", "GET"); !ok {
		zap.L().Info("策略已存在: guest, /api/v1/user/info, GET")
	}

	// exp: 将用户 "admin" 添加到 "admin" 角色
	// ps：你需要在数据库中实际拥有一个名为 "admin" 的用户和 "admin" 的角色
	if ok, _ := e.AddGroupingPolicy("admin", "admin"); !ok {
		zap.L().Info("用户角色关系已存在: admin -> admin")
	}

	// exp : 将用户 "guest" 添加到 "guest" 角色
	if ok, _ := e.AddGroupingPolicy("guest", "guest"); !ok {
		zap.L().Info("用户角色关系已存在: guest -> guest")
	}

	zap.L().Info("Casbin 初始化并加载策略完成")
	return e
}
