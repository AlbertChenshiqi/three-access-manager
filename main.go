package main

import (
	"flag"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	router2 "third-login/biz/router"
	"third-login/config"
	"third-login/pkg/redis"
)

var (
	configPath = flag.String("config", "config/config.yaml", "配置文件路径")
)

func main() {
	flag.Parse()

	// 加载配置
	err := config.LoadConfig(*configPath)
	if err != nil {
		hlog.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	initLogger(config.GlobalConfig)

	hlog.Info("Starting third-login server...")

	// 初始化Redis客户端
	redis.GetClient()
	hlog.Info("Redis client initialized")

	// 创建Hertz服务器
	h := server.Default(
		server.WithHostPorts(config.GlobalConfig.GetServerAddr()),
		server.WithReadTimeout(config.GlobalConfig.Server.ReadTimeout),
		server.WithWriteTimeout(config.GlobalConfig.Server.WriteTimeout),
	)

	// 设置路由
	router2.GeneratedRegister(h)

	// 启动服务器
	hlog.Infof("Server starting on %s", config.GlobalConfig.GetServerAddr())
	h.Spin()
}

// initLogger 初始化日志
func initLogger(cfg *config.Config) {
	// 设置日志级别
	switch cfg.Log.Level {
	case "debug":
		hlog.SetLevel(hlog.LevelDebug)
	case "info":
		hlog.SetLevel(hlog.LevelInfo)
	case "warn":
		hlog.SetLevel(hlog.LevelWarn)
	case "error":
		hlog.SetLevel(hlog.LevelError)
	default:
		hlog.SetLevel(hlog.LevelInfo)
	}

	hlog.Info("Logger initialized")
}
