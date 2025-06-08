package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kgosLj/opsvoid/internal/web/handler"
)

func RegisterRoleRouter(router *gin.RouterGroup, handler *handler.RoleHandler) {
	rg := router.Group("/role")
	{
		rg.POST("/create", handler.CreateRole)
	}
}
