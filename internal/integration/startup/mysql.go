package startup

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func InitMySQL(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("连接 Mysql 数据库失败：%s", zap.Error(err))
		panic(err)
	}
	// 设置数据库连接参数
	sql, err := db.DB()
	if err != nil {
		panic("数据库连接池设置失败")
	}
	sql.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sql.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sql.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间
	if err := sql.Ping(); err != nil {
		panic("数据库存活检查失败")
	}
	return db
}
