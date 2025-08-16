# 第三方登录管理系统

统一管理多个第三方登录应用的中间件服务，通过配置文件管理同一平台下多个应用的登录密钥，根据appid返回对应access_token。

## 功能特性

- 🔐 **统一Token管理**: 通过配置文件管理多个应用的登录密钥
- 🚀 **自动刷新机制**: 内部自动处理token获取和刷新逻辑
- 📊 **数据统计分析**: 提供access_token调用统计和监控数据
- 🔒 **应用隔离**: 通过appid区分不同应用，确保数据安全
- 📈 **扩展性强**: 支持后续添加其他第三方平台（支付宝、抖音等）
- ⚡ **高性能**: 基于Hertz框架和Redis缓存

## 技术架构

- **后端框架**: Golang + Hertz
- **API协议**: Protocol Buffers (protobuf)
- **缓存**: Redis
- **配置**: YAML配置文件
- **HTTP客户端**: 调用第三方登录API

## 项目结构

```
third-login/
├── api/
│   └── proto/           # protobuf定义文件
├── cmd/
│   └── server/          # 服务器启动入口
├── config/              # 配置文件
├── internal/
│   ├── config/          # 配置管理
│   ├── handler/         # HTTP处理器
│   └── service/         # 业务逻辑层
├── pkg/
│   ├── redis/           # Redis客户端封装
│   └── wechat/          # 微信API客户端
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 快速开始

### 环境要求

- Go 1.21+
- Redis 6.0+

### 安装依赖

```bash
make deps
```

### 配置文件

复制并修改配置文件：

```bash
cp config/config.yaml config/config.local.yaml
```

编辑 `config/config.local.yaml`，配置Redis连接和微信小程序密钥：

```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

platforms:
  wechat_miniprogram:
    apps:
      your_app_id:
        app_secret: "你的小程序密钥"
```

### 启动服务

```bash
# 开发环境
make run

# 或者指定配置文件
go run cmd/server/main.go -config=config/config.local.yaml
```

### 构建生产版本

```bash
make build-prod
```

## API接口

### 统一Token获取接口

**POST** `/api/v1/auth/token`

请求参数：
```json
{
  "platform": "wechat_miniprogram",
  "app_id": "wx1234567890abcdef"
}
```

响应：
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "access_token": "wechat_wx1234567890abcdef_1234567890abcdef"
  }
}
```
### 健康检查

**GET** `/health`

## 开发指南

### 开发环境设置

```bash
# 安装开发工具和依赖
make dev-setup

# 启动Redis（用于开发）
make redis-start

# 运行服务
make run
```

### 代码规范

```bash
# 格式化代码
make fmt

# 代码检查
make vet

# 运行测试
make test
```

### 生成protobuf文件

```bash
make http
```

## Docker部署

### 构建镜像

```bash
make docker-build
```

### 运行容器

```bash
docker run -d \
  --name third-login \
  -p 8080:8080 \
  -v $(pwd)/config:/app/config \
  third-login:latest
```

## 配置说明

### 服务器配置

```yaml
server:
  host: "0.0.0.0"        # 监听地址
  port: 8080              # 监听端口
  mode: "debug"           # 运行模式: debug/release
  read_timeout: 30s       # 读取超时
  write_timeout: 30s      # 写入超时
```

### Redis配置

```yaml
redis:
  host: "localhost"       # Redis地址
  port: 6379              # Redis端口
  password: ""            # Redis密码
  db: 0                   # 数据库编号
  pool_size: 10           # 连接池大小
```

### 平台配置

```yaml
platforms:
  wechat_miniprogram:
    name: "微信小程序"
    type: "miniprogram"
    enabled: true
    apps:
      wx1234567890abcdef:
        app_name: "测试小程序"
        app_secret: "your_app_secret"
        api_base_url: "https://api.weixin.qq.com"
        token_expire_time: 7200
        enabled: true
```

## 监控和日志

### 日志配置

```yaml
log:
  level: "info"           # 日志级别: debug/info/warn/error
  format: "json"          # 日志格式: json/text
  output: "stdout"        # 输出方式: stdout/file
```

### 统计数据

系统自动收集以下统计数据：
- 每日登录次数
- 总用户数
- 活跃token数
- 平台调用统计

## 故障排除

### 常见问题

1. **Redis连接失败**
   - 检查Redis服务是否启动
   - 验证连接配置是否正确

2. **微信API调用失败**
   - 检查app_secret是否正确
   - 验证网络连接是否正常

3. **Token过期**
   - 系统会自动刷新即将过期的token
   - 检查Redis中的token存储

### 日志查看

```bash
# 查看服务日志
tail -f logs/app.log

# 查看错误日志
grep "ERROR" logs/app.log
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请提交 Issue 或联系维护者。