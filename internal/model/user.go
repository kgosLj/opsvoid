package model

import "gorm.io/gorm"

// User 用户结构体
type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string `gorm:"not null"`
	// GORM 会自动创建一张名为 user_roles 的中间表（连接表）
	Role []*Role `gorm:"many2many:user_roles;"`
}

// Login 登录结构体

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
