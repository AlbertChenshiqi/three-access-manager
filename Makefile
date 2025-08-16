# 第三方登录管理系统 Makefile

.PHONY: build run test clean deps proto help

# 默认目标
all: build

# 构建项目
build:
	@echo "Building third-login server..."
	go build -o bin/server cmd/server/main.go
	@echo "Build completed: bin/server"

# 运行服务器
run:
	@echo "Starting third-login server..."
	go run cmd/server/main.go -config=config/config.yaml

# 运行测试
test:
	@echo "Running tests..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	rm -rf bin/
	rm -f *.log

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# 生成protobuf文件
proto:
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative api/proto/*.proto

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 代码检查
vet:
	@echo "Running go vet..."
	go vet ./...

# 安装工具
tools:
	@echo "Installing development tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 开发环境设置
dev-setup: tools deps
	@echo "Development environment setup completed"

# 生产构建
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o bin/server cmd/server/main.go

# Docker构建
docker-build:
	@echo "Building Docker image..."
	docker build -t third-login:latest .

# 启动Redis（用于开发）
redis-start:
	@echo "Starting Redis server..."
	docker run -d --name third-login-redis -p 6379:6379 redis:alpine

# 停止Redis
redis-stop:
	@echo "Stopping Redis server..."
	docker stop third-login-redis
	docker rm third-login-redis

# 检查代码质量
lint:
	@echo "Running linter..."
	golangci-lint run

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  build      - Build the server binary"
	@echo "  run        - Run the server"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build files"
	@echo "  deps       - Install dependencies"
	@echo "  proto      - Generate protobuf files"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  tools      - Install development tools"
	@echo "  dev-setup  - Setup development environment"
	@echo "  build-prod - Build for production"
	@echo "  docker-build - Build Docker image"
	@echo "  redis-start - Start Redis for development"
	@echo "  redis-stop  - Stop Redis"
	@echo "  lint       - Run linter"
	@echo "  help       - Show this help message"