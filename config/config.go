package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var config *Config

func LoadConfig(filepath string) *Config {
	v := viper.New()
	v.AddConfigPath(filepath)
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		zap.L().Error("读取 config.yaml 配置文件失败：%s", zap.Error(err))
		panic(err)
	}
	if err := v.Unmarshal(config); err != nil {
		zap.L().Error("解析 config.yaml 配置文件失败：%s", zap.Error(err))
		panic(err)
	}
	return config
}
