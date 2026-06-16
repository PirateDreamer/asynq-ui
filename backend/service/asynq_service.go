package service

import (
	"fmt"

	"github.com/hibiken/asynq"
)

// AsynqService 封装 Asynq 客户端和检查器的操作
type AsynqService struct {
	client    *asynq.Client
	inspector *asynq.Inspector
}

// NewAsynqService 创建 AsynqService 实例
func NewAsynqService(opt asynq.RedisClientOpt) *AsynqService {
	return &AsynqService{
		client:    asynq.NewClient(opt),
		inspector: asynq.NewInspector(opt),
	}
}

// InitDefaultQueue 初始化默认队列，确保面板启动就能显示
func (s *AsynqService) InitDefaultQueue() {
	task := asynq.NewTask("init", nil)
	_, _ = s.client.Enqueue(task, asynq.Queue("default"))
}

// GetQueues 获取所有队列及待处理任务数
func (s *AsynqService) GetQueues() (map[string]int, error) {
	queues, err := s.inspector.Queues()
	if err != nil {
		return nil, fmt.Errorf("获取队列列表失败: %w", err)
	}

	data := make(map[string]int, len(queues))
	for _, q := range queues {
		info, err := s.inspector.GetQueueInfo(q)
		if err != nil {
			continue
		}
		data[q] = info.Pending
	}

	return data, nil
}

// TaskGroup 任务分组结果
type TaskGroup struct {
	Pending   []map[string]interface{} `json:"pending"`
	Scheduled []map[string]interface{} `json:"scheduled"`
	Retry     []map[string]interface{} `json:"retry"`
}

// GetQueueTasks 获取队列中的任务（pending/scheduled/retry）
func (s *AsynqService) GetQueueTasks(queue string, page, pageSize int) (*TaskGroup, error) {
	opts := []asynq.ListOption{asynq.PageSize(pageSize), asynq.Page(page - 1)}

	pendingTasks, err := s.inspector.ListPendingTasks(queue, opts...)
	if err != nil {
		return nil, fmt.Errorf("获取待处理任务失败: %w", err)
	}

	scheduledTasks, err := s.inspector.ListScheduledTasks(queue, opts...)
	if err != nil {
		return nil, fmt.Errorf("获取计划任务失败: %w", err)
	}

	retryTasks, err := s.inspector.ListRetryTasks(queue, opts...)
	if err != nil {
		return nil, fmt.Errorf("获取重试任务失败: %w", err)
	}

	return &TaskGroup{
		Pending:   formatTaskInfo(pendingTasks),
		Scheduled: formatTaskInfo(scheduledTasks),
		Retry:     formatTaskInfo(retryTasks),
	}, nil
}

// PushTask 推送新任务到指定队列
func (s *AsynqService) PushTask(queue, taskType, payload string) (map[string]interface{}, error) {
	task := asynq.NewTask(taskType, []byte(payload))
	info, err := s.client.Enqueue(task, asynq.Queue(queue))
	if err != nil {
		return nil, fmt.Errorf("推送任务失败: %w", err)
	}

	return map[string]interface{}{
		"ID":      info.ID,
		"Queue":   info.Queue,
		"Type":    info.Type,
		"Payload": string(info.Payload),
		"State":   info.State,
	}, nil
}

// DeleteQueue 删除队列中的所有任务
func (s *AsynqService) DeleteQueue(name string) error {
	// 检查队列是否存在
	queues, err := s.inspector.Queues()
	if err != nil {
		return fmt.Errorf("获取队列列表失败: %w", err)
	}

	exists := false
	for _, q := range queues {
		if q == name {
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf("队列 %s 不存在", name)
	}

	// 依次删除 pending、scheduled、retry 任务
	if _, err := s.inspector.DeleteAllPendingTasks(name); err != nil {
		return fmt.Errorf("删除待处理任务失败: %w", err)
	}
	if _, err := s.inspector.DeleteAllScheduledTasks(name); err != nil {
		return fmt.Errorf("删除计划任务失败: %w", err)
	}
	if _, err := s.inspector.DeleteAllRetryTasks(name); err != nil {
		return fmt.Errorf("删除重试任务失败: %w", err)
	}

	return nil
}

// GetConsumers 获取所有消费者（运行中的 Server 实例及活跃 Worker）
func (s *AsynqService) GetConsumers() ([]map[string]interface{}, error) {
	servers, err := s.inspector.Servers()
	if err != nil {
		return nil, fmt.Errorf("获取消费者列表失败: %w", err)
	}

	result := make([]map[string]interface{}, 0, len(servers))
	for _, srv := range servers {
		workers := make([]map[string]interface{}, 0, len(srv.ActiveWorkers))
		for _, w := range srv.ActiveWorkers {
			workers = append(workers, map[string]interface{}{
				"task_id":      w.TaskID,
				"task_type":    w.TaskType,
				"task_payload": string(w.TaskPayload),
				"queue":        w.Queue,
				"started":      w.Started,
				"deadline":     w.Deadline,
			})
		}

		result = append(result, map[string]interface{}{
			"id":              srv.ID,
			"host":            srv.Host,
			"pid":             srv.PID,
			"concurrency":     srv.Concurrency,
			"queues":          srv.Queues,
			"strict_priority": srv.StrictPriority,
			"started":         srv.Started,
			"status":          srv.Status,
			"active_workers":  workers,
		})
	}

	return result, nil
}

// Close 关闭客户端连接
func (s *AsynqService) Close() {
	s.client.Close()
	s.inspector.Close()
}

// formatTaskInfo 格式化任务信息列表
func formatTaskInfo(tasks []*asynq.TaskInfo) []map[string]interface{} {
	taskList := make([]map[string]interface{}, 0, len(tasks))
	for _, t := range tasks {
		taskList = append(taskList, map[string]interface{}{
			"id":              t.ID,
			"type":            t.Type,
			"payload":         string(t.Payload),
			"queue":           t.Queue,
			"state":           t.State.String(),
			"max_retry":       t.MaxRetry,
			"retried":         t.Retried,
			"last_error":      t.LastErr,
			"timeout":         t.Timeout,
			"deadline":        t.Deadline,
			"next_process_at": t.NextProcessAt,
		})
	}
	return taskList
}
