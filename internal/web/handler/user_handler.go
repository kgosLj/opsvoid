package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/internal/service"
	"github.com/kgosLj/opsvoid/pkg/utils"
	"net/http"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

// Login 登录功能
func (h *UserHandler) Login(c *gin.Context) {
	var loginRequest model.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "请求体参数错误："+err.Error())
		return
	}

	response, err := h.svc.Login(loginRequest)
	if err == service.ErrUserNotFound {
		utils.RespondError(c, http.StatusForbidden, "登录失败："+err.Error())
		return
	} else if err != nil {
		utils.RespondError(c, http.StatusForbidden, "登录失败："+err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, response)
}

// GetUserInfo 获取用户信息(single)
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	username, exist := c.Get("username")
	if !exist {
		utils.RespondError(c, http.StatusForbidden, "无法获取当前用户权限关联信息")
		return
	}
	user, err := h.svc.GetUserInfo(username.(string))
	if err != nil {
		utils.RespondError(c, http.StatusForbidden, fmt.Sprintf("无法获取用户信息：%s", err))
		return
	}
	utils.RespondSuccess(c, http.StatusOK, user)
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var createUserRequest model.CreateUserRequest
	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "请求体参数错误："+err.Error())
		return
	}
	resp, err := h.svc.CreateUser(&createUserRequest)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "创建用户失败："+err.Error())
		return
	}
	utils.RespondSuccess(c, http.StatusOK, resp)
}

// BindRole 绑定角色
func (h *UserHandler) BindRole(c *gin.Context) {
	var bindRoleRequest model.BindRoleRequest
	if err := c.ShouldBindJSON(&bindRoleRequest); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "请求体参数错误："+err.Error())
		return
	}
	resp, err := h.svc.BindRole(&bindRoleRequest)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "绑定角色失败："+err.Error())
		return
	}
	utils.RespondSuccess(c, http.StatusOK, resp)
}
