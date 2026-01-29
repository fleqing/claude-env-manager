# 版本信息
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# 构建标志
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -s -w"

.PHONY: build run clean install test install-bin help build-linux build-mac build-all build-release

help:
	@echo "环境变量管理工具 (Go 版本) - 可用命令："
	@echo "  make build         - 编译生成可执行文件"
	@echo "  make build-release - 编译所有平台的 Release 版本（带版本信息）"
	@echo "  make run           - 运行程序"
	@echo "  make install       - 安装依赖"
	@echo "  make install-bin   - 将可执行文件安装到 /usr/local/bin (需要 sudo)"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理构建文件"
	@echo ""
	@echo "当前版本: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"

build:
	@echo "正在编译..."
	go build $(LDFLAGS) -o bin/claude-env-manager cmd/claude-env-manager/main.go
	@echo "✓ 编译完成: bin/claude-env-manager"
	@echo "版本: $(VERSION)"

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

# 交叉编译（用于 Release）
build-release:
	@echo "正在编译所有平台的 Release 版本..."
	@mkdir -p bin
	# Linux amd64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/claude-env-manager-linux-amd64 cmd/claude-env-manager/main.go
	# Linux arm64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/claude-env-manager-linux-arm64 cmd/claude-env-manager/main.go
	# macOS amd64 (Intel)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/claude-env-manager-darwin-amd64 cmd/claude-env-manager/main.go
	# macOS arm64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/claude-env-manager-darwin-arm64 cmd/claude-env-manager/main.go
	# Windows amd64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/claude-env-manager-windows-amd64.exe cmd/claude-env-manager/main.go
	@echo "✓ 所有平台编译完成 (版本: $(VERSION))"
	@ls -lh bin/

# 兼容旧命令
build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/claude-env-manager-linux cmd/claude-env-manager/main.go

build-mac:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/claude-env-manager-mac cmd/claude-env-manager/main.go

build-all: build build-linux build-mac
	@echo "✓ 所有平台编译完成"
