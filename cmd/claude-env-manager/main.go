package main

import (
	"claude-env-manager/internal/config"
	"claude-env-manager/internal/manager"
	"claude-env-manager/internal/ui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
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
