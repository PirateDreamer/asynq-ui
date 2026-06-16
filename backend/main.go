package main

import (
	"flag"
	"log"

	"backend/config"
	"backend/router"
	"backend/service"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "./config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化 Asynq 服务
	svc := service.NewAsynqService(cfg.RedisClientOpt())
	defer svc.Close()

	// 初始化默认队列
	svc.InitDefaultQueue()

	// 设置路由并启动服务
	r := router.Setup(svc)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
