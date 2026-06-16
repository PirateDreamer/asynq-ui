package router

import (
	"backend/handler"
	"backend/middleware"
	"backend/service"

	"github.com/gin-gonic/gin"
)

// Setup 初始化路由
func Setup(svc *service.AsynqService) *gin.Engine {
	r := gin.Default()

	// 跨域中间件
	r.Use(middleware.Cors())

	// 初始化 Handler
	queueHandler := handler.NewQueueHandler(svc)
	taskHandler := handler.NewTaskHandler(svc)
	consumerHandler := handler.NewConsumerHandler(svc)

	// 路由注册
	r.GET("/queues", queueHandler.GetQueues)
	r.GET("/queues/:name/tasks", taskHandler.GetQueueTasks)
	r.POST("/push", taskHandler.PushTask)
	r.DELETE("/queues/:name", queueHandler.DeleteQueue)
	r.GET("/consumers", consumerHandler.GetConsumers)

	return r
}
