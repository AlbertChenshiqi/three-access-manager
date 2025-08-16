package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	Server    *ServerConfig              `yaml:"server"`
	Redis     *RedisConfig               `yaml:"redis"`
	Platforms map[string]*PlatformConfig `yaml:"platforms"`
	Log       *LogConfig                 `yaml:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Mode         string        `yaml:"mode"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// PlatformConfig 平台配置
type PlatformConfig struct {
	Name       string                `yaml:"name"`
	Type       string                `yaml:"type"`
	Enabled    bool                  `yaml:"enabled"`
	APIBaseURL string                `yaml:"api_base_url"`
	Apps       map[string]*AppConfig `yaml:"apps"`
}

// AppConfig 应用配置
type AppConfig struct {
	AppSecret string `yaml:"app_secret"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

var GlobalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err = yaml.Unmarshal(data, &GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// GetAppConfig 获取指定平台和应用的配置
func GetAppConfig(platform, appID string) (*AppConfig, error) {
	platformConfig, exists := GlobalConfig.Platforms[platform]
	if !exists {
		return nil, fmt.Errorf("platform %s not found", platform)
	}

	if !platformConfig.Enabled {
		return nil, fmt.Errorf("platform %s is disabled", platform)
	}

	appConfig, exists := platformConfig.Apps[appID]
	if !exists {
		return nil, fmt.Errorf("app %s not found in platform %s", appID, platform)
	}

	return appConfig, nil
}

// GetRedisAddr 获取Redis连接地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetServerAddr 获取服务器监听地址
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
