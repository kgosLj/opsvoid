package startup

import (
	"github.com/kgosLj/opsvoid/internal/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitAdminUser(db *gorm.DB) {
	// 初始化管理员用户
	adminRole := &model.Role{
		Name: "admin",
		Desc: "管理员",
	}
	if err := db.FirstOrCreate(&adminRole, model.Role{Name: "admin"}).Error; err != nil {
		zap.L().Fatal("创建或查找 [admin] 角色失败", zap.Error(err))
	}

	guestRole := &model.Role{
		Name: "guest",
		Desc: "游客",
	}
	if err := db.FirstOrCreate(&guestRole, model.Role{Name: "guest"}).Error; err != nil {
		zap.L().Fatal("创建或查找 [guest] 角色失败", zap.Error(err))
	}
	zap.L().Info("角色初始化或检查完成")

	// --- 处理 Admin 用户 ---
	var adminExists bool
	if err := db.Model(&model.User{}).Where("username = ?", "admin").First(&model.User{}).Error; err == nil {
		adminExists = true
	} else if err != gorm.ErrRecordNotFound {
		zap.L().Fatal("查询 [admin] 用户失败", zap.Error(err))
	}

	if !adminExists {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			zap.L().Fatal("为 [admin] 用户哈希密码失败", zap.Error(err))
		}

		adminUser := model.User{
			Username: "admin",
			Password: string(hashedPassword),
			Role:     []*model.Role{adminRole}, // 明确关联已存在的 adminRole
		}
		if err := db.Create(&adminUser).Error; err != nil {
			zap.L().Fatal("创建 [admin] 用户失败", zap.Error(err))
		}
		zap.L().Info("初始用户 [admin] 创建成功")
	}

	// --- 处理 Guest 用户 ---
	var guestExists bool
	if err := db.Model(&model.User{}).Where("username = ?", "guest").First(&model.User{}).Error; err == nil {
		guestExists = true
	} else if err != gorm.ErrRecordNotFound {
		zap.L().Fatal("查询 [guest] 用户失败", zap.Error(err))
	}

	if !guestExists {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			zap.L().Fatal("为 [guest] 用户哈希密码失败", zap.Error(err))
		}

		guestUser := model.User{
			Username: "guest",
			Password: string(hashedPassword),
			Role:     []*model.Role{guestRole}, // 明确关联已存在的 guestRole
		}
		if err := db.Create(&guestUser).Error; err != nil {
			zap.L().Fatal("创建 [guest] 用户失败", zap.Error(err))
		}
		zap.L().Info("初始用户 [guest] 创建成功")
	}
}
