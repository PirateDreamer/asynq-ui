package handler

import (
	"net/http"

	"backend/service"

	"github.com/gin-gonic/gin"
)

// QueueHandler 队列相关处理器
type QueueHandler struct {
	svc *service.AsynqService
}

// NewQueueHandler 创建 QueueHandler
func NewQueueHandler(svc *service.AsynqService) *QueueHandler {
	return &QueueHandler{svc: svc}
}

// GetQueues 获取所有队列状态
func (h *QueueHandler) GetQueues(c *gin.Context) {
	data, err := h.svc.GetQueues()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// DeleteQueue 删除队列
func (h *QueueHandler) DeleteQueue(c *gin.Context) {
	queueName := c.Param("name")

	if err := h.svc.DeleteQueue(queueName); err != nil {
		if err.Error() == "队列 "+queueName+" 不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "queue deleted successfully"})
}
