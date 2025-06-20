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
		ug.GET("/info", handler.GetUserInfo)
		ug.POST("/create", handler.CreateUser)
		ug.POST("/bindrole", handler.BindRole)
		// TODO：addrole 功能，给用户在原有的角色基础增加权限
	}
}
