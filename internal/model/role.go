package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name  string  `gorm:"unique" exp:"admin、user、guest"`
	Desc  string  `exp:"管理员、普通用户、游客"`
	Users []*User `gorm:"many2many:user_roles;"` // 关联用户，多对多
}

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}
