package utils

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kgosLj/opsvoid/config"
	"net/http"
)

func GetDSN(db config.MysqlConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.DBName)
}

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
