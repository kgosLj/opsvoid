package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kgosLj/opsvoid/internal/web/handler"
)

// RegisterUserRouter 注册用户路由
func RegisterUserRouter(router *gin.RouterGroup, handler *handler.UserHandler) {
	ug := router.Group("/user")
	{
		ug.POST("/login", handler.Login)
		ug.POST("/create", handler.CreateUser)
	}

}
