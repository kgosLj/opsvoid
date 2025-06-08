package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadConfig(filepath string) *Config {
	v := viper.New()
	v.AddConfigPath(filepath)
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		zap.L().Error("读取 config.yaml 配置文件失败：%s", zap.Error(err))
		panic(err)
	}
	
}
