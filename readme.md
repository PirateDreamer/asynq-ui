# Asynq UI

一个轻量级的 [Asynq](https://github.com/hibiken/asynq) 任务队列可视化管理面板，提供队列监控、任务管理和消费者状态查看等功能。

## 功能特性

### 队列管理
- 查看所有队列及其待处理任务数量
- 支持删除指定队列中的所有任务（Pending / Scheduled / Retry）
- 每 5 秒自动刷新队列状态

### 任务查看
- 按队列浏览 Pending（待处理）、Scheduled（计划中）、Retry（重试中）三类任务
- 分页加载，支持自定义每页数量
- 展示任务详情：ID、类型、队列、状态、Payload、下次执行时间、错误信息等

### 推送任务
- 手动向指定队列推送新任务
- 内置 JSON 编辑器（CodeMirror），支持语法高亮、实时校验、格式化/压缩

### 消费者监控
- 查看所有 Asynq Server 实例（主机、PID、并发数、监听队列）
- 查看各实例的活跃 Worker 列表及其正在执行的任务详情
- 实时状态标识（Active / Stopped）

## 技术栈

| 组件 | 技术 |
|------|------|
| 后端 | Go + [Gin](https://github.com/gin-gonic/gin) + [Asynq](https://github.com/hibiken/asynq) |
| 前端 | 原生 HTML/CSS/JS + [CodeMirror](https://codemirror.net/) |
| 配置 | [Viper](https://github.com/spf13/viper)（支持热加载） |
| 存储 | Redis |

## 项目结构

```
asynq-ui/
├── backend/
│   ├── main.go                # 入口：解析参数、组装模块、启动服务
│   ├── config/
│   │   └── config.go          # 配置加载（Viper + 热加载）
│   ├── service/
│   │   └── asynq_service.go   # Asynq 客户端封装（队列/任务/消费者操作）
│   ├── handler/
│   │   ├── queue_handler.go   # 队列相关 Handler
│   │   ├── task_handler.go    # 任务相关 Handler
│   │   └── consumer_handler.go# 消费者相关 Handler
│   ├── middleware/
│   │   └── cors.go            # CORS 跨域中间件
│   ├── router/
│   │   └── router.go          # 路由注册
│   ├── config.yaml            # Redis 配置（gitignore，不提交）
│   ├── config-example.yaml    # Redis 配置示例
│   ├── go.mod
│   └── go.sum
├── frontend/
│   └── index.html             # 前端页面（单文件，无构建步骤）
└── readme.md
```

## 快速开始

### 前置条件

- Go >= 1.21
- Redis 已启动并可访问

### 1. 配置 Redis 连接

```bash
cp backend/config-example.yaml backend/config.yaml
```

编辑 `backend/config.yaml`：

```yaml
redis:
  host: "127.0.0.1"
  port: "6379"
  password: ""
  index: 0
```

### 2. 启动后端

```bash
go run ./backend/
```

服务默认运行在 `http://localhost:8080`。可通过 `-config` 参数指定配置文件路径：

```bash
go run ./backend/ -config /path/to/config.yaml
```

### 3. 打开前端

直接用浏览器打开 `frontend/index.html`，或使用 VS Code Live Server 等工具启动。

> 前端默认请求 `http://localhost:8080`，请确保后端已启动。

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/queues` | 获取所有队列及待处理任务数 |
| `GET` | `/queues/:name/tasks` | 获取队列任务列表（支持 `page`、`pageSize` 分页参数） |
| `POST` | `/push` | 推送新任务（Body: `{queue, type, payload}`） |
| `DELETE` | `/queues/:name` | 删除队列中的所有任务 |
| `GET` | `/consumers` | 获取所有消费者实例及活跃 Worker |

## 界面预览

**队列状态面板** — 展示所有队列及任务数量，支持查看任务和清空队列：

```
┌──────────────────────────────────────────────────┐
│  Queue Status                                    │
├──────────┬──────────┬────────────────────────────┤
│  队列名   │  待处理数  │  操作                       │
├──────────┼──────────┼────────────────────────────┤
│  default │  42      │  [查看任务]  [删除队列]       │
│  email   │  15      │  [查看任务]  [删除队列]       │
└──────────┴──────────┴────────────────────────────┘
```

**消费者面板** — 展示所有 Server 实例及 Worker 状态：

```
┌───────────────────────────────────────────────────────────────┐
│  Consumers                                                    │
├─────────────┬─────┬────────┬──────────┬────────┬──────────────┤
│  Host/PID   │ 并发 │ 队列    │ 活跃Worker│ 状态   │  启动时间     │
├─────────────┼─────┼────────┼──────────┼────────┼──────────────┤
│  host:12345 │  10 │ default│  3       │ Active │  2026-01-01  │
└─────────────┴─────┴────────┴──────────┴────────┴──────────────┘
```

## License

[MIT](LICENSE)
