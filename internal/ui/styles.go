package ui

import "github.com/charmbracelet/lipgloss"

var (
	// 颜色定义
	primaryColor   = lipgloss.Color("#7D56F4")
	successColor   = lipgloss.Color("#04B575")
	errorColor     = lipgloss.Color("#FF4672")
	warningColor   = lipgloss.Color("#FFB86C")
	subtleColor    = lipgloss.Color("#6C7086")
	activeColor    = lipgloss.Color("#A6E3A1")
	inactiveColor  = lipgloss.Color("#585B70")

	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// 菜单项样式
	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	// 选中的菜单项样式
	selectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(primaryColor).
				Bold(true)

	// 成功消息样式
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// 错误消息样式
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// 警告消息样式
	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// 提示文本样式
	subtleStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// 表格样式
	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true)

	tableCellStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	// 激活状态样式
	activeStatusStyle = lipgloss.NewStyle().
				Foreground(activeColor).
				Bold(true)

	// 未激活状态样式
	inactiveStatusStyle = lipgloss.NewStyle().
				Foreground(inactiveColor)

	// 输入框样式
	inputStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// 帮助文本样式
	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			MarginTop(1)
)
