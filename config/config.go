package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Allconfig *Config

func LoadConfig(filepath string) *Config {
	Allconfig = &Config{}
	v := viper.New()
	v.SetConfigFile(filepath)
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		zap.L().Error("读取 config.yaml 配置文件失败：%s", zap.Error(err))
		panic(err)
	}

	if err := v.Unmarshal(Allconfig); err != nil {
		zap.L().Error("解析 config.yaml 配置文件失败：%s", zap.Error(err))
		panic(err)
	}
	return Allconfig
}
