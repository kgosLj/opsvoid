package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kgosLj/opsvoid/internal/web/handler"
)

func RegisterRoleRouter(router *gin.RouterGroup, handler *handler.RoleHandler) {
	rg := router.Group("/role")
	{
		// TODO： fix: 创建role的时候需要创建 casbin 角色规则，并且指定路径访问的权限
		rg.POST("/createrole", handler.CreateRole)
		rg.POST("/createrbac", handler.CreateRbac)
		// TODO: Update role 功能
		// TODO: 1. 获取全部 role 功能
		// TODO: 2. 删除指定 role 功能
		// TODO: 3. 添加 role 权限能访问的路径功能
		// TODO: 4. 删除 role 权限访问的功能
	}
}
