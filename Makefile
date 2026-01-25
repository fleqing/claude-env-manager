.PHONY: build run clean install test install-bin help

help:
	@echo "环境变量管理工具 (Go 版本) - 可用命令："
	@echo "  make build       - 编译生成可执行文件"
	@echo "  make run         - 运行程序"
	@echo "  make install     - 安装依赖"
	@echo "  make install-bin - 将可执行文件安装到 /usr/local/bin (需要 sudo)"
	@echo "  make test        - 运行测试"
	@echo "  make clean       - 清理构建文件"

build:
	@echo "正在编译..."
	go build -o bin/claude-env-manager cmd/claude-env-manager/main.go
	@echo "✓ 编译完成: bin/claude-env-manager"

run:
	go run cmd/claude-env-manager/main.go

install:
	go mod download
	go mod tidy

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

install-bin:
	@echo "正在安装 claude-env-manager 到 /usr/local/bin..."
	@if [ ! -f bin/claude-env-manager ]; then \
		echo "错误: bin/claude-env-manager 不存在，请先运行 'make build'"; \
		exit 1; \
	fi
	sudo cp -f bin/claude-env-manager /usr/local/bin/claude-env-manager
	@echo "✓ 安装完成: /usr/local/bin/claude-env-manager"
	@echo "现在可以直接使用 'claude-env-manager' 命令了"

# 交叉编译
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/claude-env-manager-linux cmd/claude-env-manager/main.go

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/claude-env-manager-mac cmd/claude-env-manager/main.go

build-all: build build-linux build-mac
	@echo "✓ 所有平台编译完成"
