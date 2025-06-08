package main

import (
	"flag"
	"github.com/kgosLj/opsvoid/config"
	"github.com/kgosLj/opsvoid/internal/integration"
)

var configPath string

func main() {
	// 指定配置文件路径
	flag.StringVar(&configPath, "config", "./configs/config.yaml", "配置文件路径")
	flag.Parse()

	// 读取配置文件
	config := config.LoadConfig(configPath)

	// 初始化全局的方法
	integration.Initizalizate(config)
}
