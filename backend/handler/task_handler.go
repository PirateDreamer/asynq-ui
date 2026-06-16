package handler

import (
	"net/http"
	"strconv"

	"backend/service"

	"github.com/gin-gonic/gin"
)

// TaskHandler 任务相关处理器
type TaskHandler struct {
	svc *service.AsynqService
}

// NewTaskHandler 创建 TaskHandler
func NewTaskHandler(svc *service.AsynqService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// GetQueueTasks 获取队列中的任务
func (h *TaskHandler) GetQueueTasks(c *gin.Context) {
	queueName := c.Param("name")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "30"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 30
	}

	result, err := h.svc.GetQueueTasks(queueName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PushTask 推送新任务
func (h *TaskHandler) PushTask(c *gin.Context) {
	var req struct {
		Queue   string `json:"queue"`
		Type    string `json:"type"`
		Payload string `json:"payload"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.PushTask(req.Queue, req.Type, req.Payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
