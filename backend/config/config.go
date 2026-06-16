package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

// AppConfig 应用配置
type AppConfig struct {
	Redis RedisConfig
}

// RedisConfig Redis 连接配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Index    int
}

// Load 加载配置文件
func Load(path string) (*AppConfig, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 热加载
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("\n【检测到配置变动】文件: %s, 操作: %s\n", e.Name, e.Op)
	})
	viper.WatchConfig()

	cfg := &AppConfig{
		Redis: RedisConfig{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetString("redis.port"),
			Password: viper.GetString("redis.password"),
			Index:    viper.GetInt("redis.index"),
		},
	}

	return cfg, nil
}

// RedisClientOpt 返回 Asynq Redis 连接配置
func (c *AppConfig) RedisClientOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.Index,
	}
}
