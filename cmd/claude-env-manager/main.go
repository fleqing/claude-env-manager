package main

import (
	"claude-env-manager/internal/config"
	"claude-env-manager/internal/manager"
	"claude-env-manager/internal/ui"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	// 版本信息，通过 ldflags 在编译时注入
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	// 添加版本参数
	version := flag.Bool("version", false, "显示版本信息")
	v := flag.Bool("v", false, "显示版本信息（简写）")
	flag.Parse()

	if *version || *v {
		fmt.Printf("claude-env-manager %s\n", Version)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build Time: %s\n", BuildTime)
		os.Exit(0)
	}
	// 创建配置
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法创建配置: %v\n", err)
		os.Exit(1)
	}

	// 创建管理器
	mgr, err := manager.NewManager(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法初始化管理器: %v\n", err)
		os.Exit(1)
	}

	// 创建 UI 模型
	model := ui.NewModel(mgr)

	// 启动程序
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
