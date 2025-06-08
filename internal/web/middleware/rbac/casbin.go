package rbac

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kgosLj/opsvoid/pkg/utils"
)

func CasbinMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在这里添加忽略的路径
		if c.Request.URL.Path == "/api/v1/user/login" {
			c.Next()
			return
		}

		username, exist := c.Get("username")
		if !exist {
			utils.RespondError(c, 403, "无法获取当前用户权限关联信息")
			c.Abort()
			return
		}
		// 将 username 转换为字符串类型
		sub, ok := username.(string)
		if !ok {
			utils.RespondError(c, 500, "上下文中存储的用户信息格式不正确")
			c.Abort()
			return
		}
		// 获取请求的URI
		obj := c.Request.URL.Path
		// 获取请求的方法
		act := c.Request.Method
		// 执行策略匹配
		ok, err := e.Enforce(sub, obj, act)
		if err != nil {
			utils.RespondError(c, 500, "内部错误")
			c.Abort()
			return
		}
		if ok {
			c.Next()
		} else {
			utils.RespondError(c, 403, "权限不足, 禁止访问（请联系管理员开通权限）")
			c.Abort()
			return
		}
	}
}
