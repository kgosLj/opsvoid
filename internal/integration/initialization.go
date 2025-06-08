package integration

import (
	"github.com/kgosLj/opsvoid/config"
	"github.com/kgosLj/opsvoid/internal/integration/startup"
	"github.com/kgosLj/opsvoid/internal/repository"
	"github.com/kgosLj/opsvoid/internal/repository/dao"
	"github.com/kgosLj/opsvoid/pkg/logger"
	"github.com/kgosLj/opsvoid/pkg/utils"
	"gorm.io/gorm"
)

// Initizalizate 初始化全局函数
func Initizalizate(config *config.Config) {
	// 初始化日志
	logger.InitZapLogger()

	// 初始化数据库
	dsn := utils.GetDSN(config.Mysql)
	db := startup.InitMySQL(dsn)

	// 初始化 gin 服务
	dao := InitDao(db)
	repository := InitRepository(dao)
}

// InitDao 初始化 DAO
func InitDao(db *gorm.DB) *AppDao {
	appDao := new(AppDao)
	appDao.UserDao = dao.NewUserDao(db)
	return appDao
}

func InitRepository(dao *AppDao) *AppRepository {
	appRepository := new(AppRepository)
	appRepository.UserRepository = repository.NewUserRepository(dao.UserDao)
	return appRepository
}
