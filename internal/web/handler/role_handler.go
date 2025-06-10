package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/service"
	"github.com/kgosLj/opsvoid/pkg/utils"
	"net/http"
)

type RoleHandler struct {
	svc service.RoleService
}

func NewRoleHandler(svc service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// CreateRole 创建角色
// 创建时候可以为角色添加权限
// 创建完成之后就通过 CreateRbac 去添加角色的权限
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req model.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "请求体参数错误："+err.Error())
		return
	}
	if err := h.svc.CreateRole(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "创建角色失败："+err.Error())
		return
	}
	utils.RespondSuccess(c, http.StatusOK, "创建角色成功！")
}

func (h *RoleHandler) CreateRbac(c *gin.Context) {
	
}
