package redis

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"sync"
	"third-login/config"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client Redis客户端封装
type Client struct {
	rdb *redis.Client
}

var redisClient *Client
var redisOnce sync.Once

// NewClient 创建Redis客户端
func NewClient(cfg *config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// InitClient 初始化全局微信客户端
func InitClient() {
	redisOnce.Do(func() {
		client, err := NewClient(config.GlobalConfig.Redis)
		if err != nil {
			hlog.Fatalf("Failed to connect to Redis: %v", err)
		}
		redisClient = client
	})
}

func GetClient() *Client {
	if redisClient == nil {
		InitClient()
	}
	return redisClient
}

// Close 关闭Redis连接
func (c *Client) Close() error {
	return c.rdb.Close()
}

// SetAccessToken 存储access_token
func (c *Client) SetAccessToken(ctx context.Context, platform, appID, token string, expiration time.Duration) error {
	key := fmt.Sprintf("access_token:%s:%s", platform, appID)
	return c.rdb.Set(ctx, key, token, expiration).Err()
}

// GetAccessToken 获取access_token
func (c *Client) GetAccessToken(ctx context.Context, platform, appID string) (string, error) {
	key := fmt.Sprintf("access_token:%s:%s", platform, appID)
	return c.rdb.Get(ctx, key).Result()
}

// AccessTokenTTL access_token ttl
func (c *Client) AccessTokenTTL(ctx context.Context, platform, appID string) (time.Duration, error) {
	key := fmt.Sprintf("access_token:%s:%s", platform, appID)
	return c.rdb.TTL(ctx, key).Result()
}

// SetUserSession 存储用户会话信息
func (c *Client) SetUserSession(ctx context.Context, platform, appID, openID string, sessionData map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s:%s:%s", platform, appID, openID)
	return c.rdb.HMSet(ctx, key, sessionData).Err()
}

// GetUserSession 获取用户会话信息
func (c *Client) GetUserSession(ctx context.Context, platform, appID, openID string) (map[string]string, error) {
	key := fmt.Sprintf("session:%s:%s:%s", platform, appID, openID)
	return c.rdb.HGetAll(ctx, key).Result()
}

// DeleteUserSession 删除用户会话
func (c *Client) DeleteUserSession(ctx context.Context, platform, appID, openID string) error {
	key := fmt.Sprintf("session:%s:%s:%s", platform, appID, openID)
	return c.rdb.Del(ctx, key).Err()
}

// IncrementStats 增加统计计数
func (c *Client) IncrementStats(ctx context.Context, key string) error {
	return c.rdb.Incr(ctx, key).Err()
}

// IncrementDailyStats 增加每日统计
func (c *Client) IncrementDailyStats(ctx context.Context, platform, appID, date string) error {
	key := fmt.Sprintf("stats:login:daily:%s:%s:%s", platform, appID, date)
	return c.rdb.Incr(ctx, key).Err()
}

// IncrementTotalStats 增加总统计
func (c *Client) IncrementTotalStats(ctx context.Context, platform, appID string) error {
	key := fmt.Sprintf("stats:login:total:%s:%s", platform, appID)
	return c.rdb.Incr(ctx, key).Err()
}

// SetAppStats 设置应用统计信息
func (c *Client) SetAppStats(ctx context.Context, appID string, field string, value interface{}) error {
	key := fmt.Sprintf("stats:app:%s", appID)
	return c.rdb.HSet(ctx, key, field, value).Err()
}

// GetAppStats 获取应用统计信息
func (c *Client) GetAppStats(ctx context.Context, appID string) (map[string]string, error) {
	key := fmt.Sprintf("stats:app:%s", appID)
	return c.rdb.HGetAll(ctx, key).Result()
}

// SetPlatformStats 设置平台统计信息
func (c *Client) SetPlatformStats(ctx context.Context, platform string, field string, value interface{}) error {
	key := fmt.Sprintf("stats:platform:%s", platform)
	return c.rdb.HSet(ctx, key, field, value).Err()
}

// GetPlatformStats 获取平台统计信息
func (c *Client) GetPlatformStats(ctx context.Context, platform string) (map[string]string, error) {
	key := fmt.Sprintf("stats:platform:%s", platform)
	return c.rdb.HGetAll(ctx, key).Result()
}

// GetDailyStats 获取指定日期范围的每日统计
func (c *Client) GetDailyStats(ctx context.Context, platform, appID string, dates []string) (map[string]int64, error) {
	result := make(map[string]int64)
	for _, date := range dates {
		key := fmt.Sprintf("stats:login:daily:%s:%s:%s", platform, appID, date)
		val, err := c.rdb.Get(ctx, key).Int64()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		result[date] = val
	}
	return result, nil
}

// Exists 检查key是否存在
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.rdb.Exists(ctx, key).Result()
	return result > 0, err
}

// TTL 获取key的过期时间
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.rdb.TTL(ctx, key).Result()
}

// Expire 设置key的过期时间
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.rdb.Expire(ctx, key, expiration).Err()
}
