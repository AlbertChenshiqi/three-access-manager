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

- **后端框架**: Golang + CloudWego Hertz
- **API协议**: HTTP RESTful API + Protocol Buffers (protobuf)
- **缓存**: Redis
- **配置**: YAML配置文件
- **HTTP客户端**: 调用第三方登录API
- **代码生成**: Hertz代码生成工具

## 项目结构

```
third-login/
├── .hz                     # Hertz代码生成配置
├── api/
│   ├── api.proto           # API定义文件
│   └── http/               # HTTP接口定义
│       ├── auth/           # 认证相关接口
│       └── service.proto   # 服务定义
├── biz/                    # 业务逻辑层
│   ├── handler/            # 请求处理器
│   │   ├── http/           # HTTP处理器
│   │   └── ping.go         # 健康检查
│   ├── middleware/         # 中间件
│   ├── model/              # 数据模型
│   │   ├── api/            # API模型
│   │   └── http/           # HTTP模型
│   ├── router/             # 路由配置
│   └── service/            # 业务服务
│       └── auth.go         # 认证服务
├── bin/                    # 编译输出目录
├── cmd/                    # 命令行工具
├── config/                 # 配置文件
│   ├── config.go           # 配置结构定义
│   └── config.yaml         # 配置文件
├── pkg/                    # 公共包
│   ├── redis/              # Redis客户端封装
│   └── wechat/             # 微信API客户端
├── script/                 # 脚本文件
├── main.go                 # 主程序入口
├── build.sh                # 构建脚本
├── Makefile                # 构建配置
└── README.md
```

## 快速开始

### 环境要求

- Go 1.21+
- Redis 6.0+
- Protocol Buffers 编译器 (protoc)

### 安装依赖

```bash
# 安装Go依赖
go mod tidy

# 或使用Makefile
make deps
```

### 配置文件

编辑 `config/config.yaml`，配置Redis连接和微信小程序密钥：

```yaml
server:
  host: "0.0.0.0"
  port: 8081
  mode: "debug"
  read_timeout: 30s
  write_timeout: 30s

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10

platforms:
  wechat_miniprogram:
    name: "微信小程序"
    type: "wechat_miniprogram"
    enabled: true
    api_base_url: "https://api.weixin.qq.com"
    apps:
      wx1234567890abcdef:
        app_secret: "your_app_secret_here"
      wx9876543210fedcba:
        app_secret: "your_prod_app_secret_here"
```

### 启动服务

```bash
# 直接运行
go run main.go -config=config/config.yaml

# 或使用Makefile
make run

# 构建后运行
make build
./bin/server -config=config/config.yaml
```

### 构建生产版本

```bash
# 生产构建
make build-prod

# 或使用构建脚本
./build.sh
```

## API接口

### 健康检查

**GET** `/ping`

响应：
```json
{
  "message": "pong"
}
```

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
    "access_token": "wechat_wx1234567890abcdef_1234567890abcdef",
    "expires_in": 7200
  }
}
```

### 获取统计数据

**GET** `/api/v1/dashboard/stats`

查询参数：
- `platform`: 平台标识
- `app_id`: 应用ID
- `date_range`: 日期范围 (today, week, month)

### 获取分析数据

**GET** `/api/v1/dashboard/analytics`

查询参数：
- `platform`: 平台标识
- `app_id`: 应用ID
- `date_range`: 日期范围

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

### 代码生成

本项目使用Hertz代码生成工具，基于protobuf定义自动生成代码：

```bash
# 生成HTTP代码
hz new -module third-login

# 更新代码
hz update
```

### 代码规范

```bash
# 格式化代码
make fmt

# 代码检查
make vet

# 运行测试
make test

# 代码质量检查
make lint
```

### 生成protobuf文件

```bash
# 生成protobuf Go代码
make proto

# 或直接使用protoc
protoc --go_out=. --go_opt=paths=source_relative api/*.proto
```

## 配置说明

### 服务器配置

```yaml
server:
  host: "0.0.0.0"        # 监听地址
  port: 8081              # 监听端口
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
  min_idle_conns: 5       # 最小空闲连接数
  dial_timeout: 5s        # 连接超时
  read_timeout: 3s        # 读取超时
  write_timeout: 3s       # 写入超时
```

### 平台配置

```yaml
platforms:
  wechat_miniprogram:
    name: "微信小程序"
    type: "wechat_miniprogram"
    enabled: true
    api_base_url: "https://api.weixin.qq.com"
    apps:
      wx1234567890abcdef:
        app_secret: "your_app_secret"
```

### 日志配置

```yaml
log:
  level: "info"           # 日志级别: debug/info/warn/error
  format: "json"          # 日志格式: json/text
  output: "stdout"        # 输出方式: stdout/file
  file_path: "logs/app.log" # 日志文件路径
  max_size: 100           # 最大文件大小(MB)
  max_backups: 3          # 备份文件数量
  max_age: 28             # 保留天数
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
  -p 8081:8081 \
  -v $(pwd)/config:/app/config \
  third-login:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  third-login:
    build: .
    ports:
      - "8081:8081"
    volumes:
      - ./config:/app/config
    depends_on:
      - redis
    environment:
      - CONFIG_PATH=/app/config/config.yaml

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```

## 监控和日志

### 统计数据

系统自动收集以下统计数据：
- 每日登录次数
- 总用户数
- 活跃token数
- 平台调用统计

### 日志查看

```bash
# 查看服务日志
tail -f logs/app.log

# 查看错误日志
grep "ERROR" logs/app.log

# 实时查看日志
docker logs -f third-login
```

## 故障排除

### 常见问题

1. **Redis连接失败**
   - 检查Redis服务是否启动
   - 验证连接配置是否正确
   - 检查网络连接和防火墙设置

2. **微信API调用失败**
   - 检查app_secret是否正确
   - 验证网络连接是否正常
   - 确认微信API接口地址正确

3. **Token过期**
   - 系统会自动刷新即将过期的token
   - 检查Redis中的token存储
   - 验证配置的过期时间设置

4. **端口占用**
   - 修改配置文件中的端口号
   - 检查是否有其他服务占用端口

### 调试模式

```bash
# 启用调试模式
export LOG_LEVEL=debug
go run main.go -config=config/config.yaml
```

## 性能优化

### Redis优化

- 合理设置连接池大小
- 使用Redis集群提高可用性
- 定期清理过期数据

### 应用优化

- 启用生产模式 (`mode: release`)
- 合理设置超时时间
- 使用负载均衡

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### 开发规范

- 遵循Go代码规范
- 添加必要的单元测试
- 更新相关文档
- 确保所有测试通过

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请提交 Issue 或联系维护者。

---

**注意**: 本项目基于CloudWego Hertz框架开发，具有高性能和易扩展的特点。更多关于Hertz的信息请参考 [CloudWego官方文档](https://www.cloudwego.io/zh/docs/hertz/)。