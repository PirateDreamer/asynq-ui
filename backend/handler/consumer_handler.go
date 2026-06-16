package handler

import (
	"net/http"

	"backend/service"

	"github.com/gin-gonic/gin"
)

// ConsumerHandler 消费者相关处理器
type ConsumerHandler struct {
	svc *service.AsynqService
}

// NewConsumerHandler 创建 ConsumerHandler
func NewConsumerHandler(svc *service.AsynqService) *ConsumerHandler {
	return &ConsumerHandler{svc: svc}
}

// GetConsumers 获取所有消费者
func (h *ConsumerHandler) GetConsumers(c *gin.Context) {
	result, err := h.svc.GetConsumers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
