package utils

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kgosLj/opsvoid/config"
	"github.com/kgosLj/opsvoid/internal/integration/startup"
	"go.uber.org/zap"
	"net/http"
)

// GetDSN 获取数据库连接字符串
func GetDSN(db config.MysqlConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.DBName)
}

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有域名跨域访问
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400, // 24h
	})
}

// RespondSuccess 成功响应
func RespondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  "success",
	})
}

// RespondError 错误响应
func RespondError(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, gin.H{
		"code": httpCode,
		"msg":  msg,
	})
}

// IsRbacPolicyExists 判断 rbac 策略是否存在
func IsRbacPolicyExists(sub, obj, act string) bool {
	exist, err := startup.E.HasPolicy(sub, obj, act)
	if err != nil {
		zap.L().Info("获取 rbac 权限时候出错", zap.Error(err))
		return false
	}
	if exist {
		return true
	}
	return false
}
