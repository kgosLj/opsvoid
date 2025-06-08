package integration

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kgosLj/opsvoid/config"
	"github.com/kgosLj/opsvoid/internal/integration/startup"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/repository"
	"github.com/kgosLj/opsvoid/internal/repository/dao"
	"github.com/kgosLj/opsvoid/internal/service"
	"github.com/kgosLj/opsvoid/internal/web/handler"
	"github.com/kgosLj/opsvoid/internal/web/middleware/jwt"
	"github.com/kgosLj/opsvoid/internal/web/middleware/rbac"
	"github.com/kgosLj/opsvoid/internal/web/router"
	"github.com/kgosLj/opsvoid/pkg/logger"
	"github.com/kgosLj/opsvoid/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Initizalizate 初始化全局函数
func Initizalizate(config *config.Config) {
	// 初始化日志
	logger.InitZapLogger()

	// 初始化数据库
	dsn := utils.GetDSN(config.Mysql)
	db := startup.InitMySQL(dsn)

	// gorm 自动表迁移
	err := db.AutoMigrate(&model.User{}, &model.Role{})
	if err != nil {
		zap.L().Fatal("数据库自动迁移失败", zap.Error(err))
	}
	zap.L().Info("数据库自动迁移成功")

	// 初始化管理员用户
	startup.InitAdminUser(db)
	// 初始化 enforce
	e := startup.InitEnforce(db)

	// 初始化 gin 服务
	dao := InitDao(db)
	repository := InitRepository(dao)
	service := InitService(repository)
	handler := InitHandler(service)
	router := InitRouter(handler, e)

	// 启动服务
	port := config.Server.Port
	if port == 0 {
		port = 8080
	}
	startupPort := fmt.Sprintf(":%d", port)
	router.Run(startupPort)
}

// InitDao 初始化 DAO
func InitDao(db *gorm.DB) *AppDao {
	appDao := new(AppDao)
	appDao.UserDao = dao.NewUserDao(db)
	appDao.RoleDao = dao.NewRoleDao(db)
	return appDao
}

// InitRepository 初始化 Repository
func InitRepository(dao *AppDao) *AppRepository {
	appRepository := new(AppRepository)
	if dao.UserDao != nil {
		appRepository.UserRepository = repository.NewUserRepository(dao.UserDao)
	}
	if dao.RoleDao != nil {
		appRepository.RoleRepository = repository.NewRoleRepository(dao.RoleDao)
	}
	zap.L().Info("repository 初始化完成")

	return appRepository
}

// InitService 初始化 Service
func InitService(repository *AppRepository) *AppService {
	appService := new(AppService)
	if repository.UserRepository != nil {
		appService.UserService = service.NewUserService(repository.UserRepository)
	}
	if repository.RoleRepository != nil {
		appService.RoleService = service.NewRoleService(repository.RoleRepository)
	}
	zap.L().Info("service 初始化完成")
	return appService
}

// InitHandler 初始化 Handler
func InitHandler(service *AppService) *AppHandler {
	appHandler := new(AppHandler)
	if service.UserService != nil {
		appHandler.UserHandler = handler.NewUserHandler(service.UserService)
	}
	if service.RoleService != nil {
		appHandler.RoleHandler = handler.NewRoleHandler(service.RoleService)
	}
	return appHandler
}

// InitRouter 初始化路由
func InitRouter(handler *AppHandler, e *casbin.Enforcer) *gin.Engine {
	r := gin.Default()
	// 初始化 zap gin 路由
	r.Use(logger.GinZapMiddleware(zap.L()))
	// 处理跨域 cors 中间件
	r.Use(utils.Cors())

	apiV1 := r.Group("/api/v1")
	{
		// 在这里添加路由认证的中间件，记得忽略掉登录功能
		apiV1.Use(jwt.JWTMiddleware())
		apiV1.Use(rbac.CasbinMiddleware(e))
		if handler.UserHandler != nil {
			router.RegisterUserRouter(apiV1, handler.UserHandler)
		}
		if handler.RoleHandler != nil {
			router.RegisterRoleRouter(apiV1, handler.RoleHandler)
		}
	}
	return r

}
