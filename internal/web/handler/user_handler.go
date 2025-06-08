package handler

import (
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
