# 环境变量管理工具

一个用于管理 `~/.zshrc` 文件中 ANTHROPIC_BASE_URL 和 ANTHROPIC_AUTH_TOKEN 环境变量组合的命令行工具。

## 功能特性

- 📋 **查看组合**：列出所有环境变量组合及其激活状态
- 🔄 **快速切换**：一键切换不同的环境变量组合
- ⚡ **测速功能**：测试 API 端点的响应速度，帮助选择最快的服务
- ✏️ **编辑组合**：修改组合的名称、BASE_URL 或 AUTH_TOKEN
- ➕ **添加组合**：添加新的环境变量组合
- 🗑️ **删除组合**：删除不需要的组合
- 💾 **自动备份**：每次修改前自动备份，保留最近 10 个备份
- 🎨 **友好界面**：使用 Bubble Tea 提供美观的交互式终端界面

## 安装

### 方式一：从源码编译

1. 克隆或下载此项目
2. 安装依赖：

```bash
make install
```

3. 编译程序：

```bash
make build
```

4. （可选）安装到系统路径：

```bash
make install-bin
```

安装后可以在任何位置直接使用 `claude-env-manager` 命令。

### 方式二：直接运行

```bash
make run
```

或直接运行：

```bash
go run cmd/claude-env-manager/main.go
```

## 使用方法

如果已安装到系统路径：

```bash
claude-env-manager
```

或使用 make 命令：

```bash
make run
```

### 主菜单选项

1. **切换环境变量组合**：选择并激活某个组合
2. **测速**：测试所有 API 端点的响应速度
3. **编辑组合**：修改组合的名称、BASE_URL 或 AUTH_TOKEN
4. **添加新组合**：添加新的环境变量组合
5. **删除组合**：删除指定的组合
6. **退出**：退出程序

### 使用提示

- 切换或修改环境变量后，需要执行 `source ~/.zshrc` 或重启终端使更改生效
- 所有修改操作都会自动创建备份文件
- 备份文件保存在 `~/.claude-env-manager/backups/` 目录
- 使用方向键（↑/↓）或 vim 风格的 j/k 键导航菜单
- 按 Enter 键选择，按 q 键退出

## 测速功能

测速功能会并发测试所有配置的 API 端点，显示：

- 响应时间（毫秒）
- 连接状态（成功/失败）
- 实时进度显示

帮助您快速找到响应最快的 API 服务。

## 环境变量格式

工具识别以下格式的环境变量组合：

```bash
# 组合名称
export ANTHROPIC_BASE_URL=https://api.example.com
export ANTHROPIC_AUTH_TOKEN=your_token_here
```

停用的组合会被注释：

```bash
# 组合名称
#export ANTHROPIC_BASE_URL=https://api.example.com
#export ANTHROPIC_AUTH_TOKEN=your_token_here
```

## 项目结构

```
env-manager/
├── cmd/
│   └── claude-env-manager/
│       └── main.go          # 程序入口
├── internal/
│   ├── config/              # 配置管理
│   ├── manager/             # 环境变量管理器
│   ├── model/               # 数据模型
│   ├── parser/              # 配置文件解析器
│   ├── speedtest/           # 测速功能
│   └── ui/                  # 用户界面
├── bin/                     # 编译输出目录
├── go.mod                   # Go 模块定义
├── Makefile                 # 构建脚本
└── README.md                # 项目文档
```

## 技术栈

- Go 1.25+
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - 终端 UI 框架
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI 组件库
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - 样式和布局

## 开发

### 运行测试

```bash
make test
```

### 清理构建文件

```bash
make clean
```

### 交叉编译

```bash
# 编译 Linux 版本
make build-linux

# 编译 macOS 版本
make build-mac

# 编译所有平台
make build-all
```

## 许可

MIT License
