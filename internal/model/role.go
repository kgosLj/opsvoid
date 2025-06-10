package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name  string  `gorm:"unique" exp:"admin、user、guest"`
	Desc  string  `exp:"管理员、普通用户、游客"`
	Users []*User `gorm:"many2many:user_roles;"` // 关联用户，多对多
}

// RoleCreateRequest 创建角色请求 (包含 rbac 策略的添加 -- 必选项)
type RoleCreateRequest struct {
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc" validate:"required"`
	// Sub  string `json:"sub" validate:"required"` 后端默认用 Name 作为 Sub 对象
	Obj string `json:"obj" validate:"required"`
	Act string `json:"act" validate:"required"`
}
