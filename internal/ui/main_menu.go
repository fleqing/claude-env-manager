package ui

import (
	"claude-env-manager/internal/speedtest"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// 主菜单操作选项（不包括组合列表，组合列表动态显示）
var mainMenuOptions = []string{
	"切换环境变量组合",
	"测速",
	"编辑组合",
	"添加新组合",
	"删除组合",
	"退出",
}

// updateMainMenu 更新主菜单
func (m Model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	groups := m.manager.GetGroups()
	separatorIndex := len(groups)
	totalItems := len(groups) + 1 + len(mainMenuOptions) // 组合 + 分隔线 + 操作选项

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.cursor--
			// 光标最小值为第一个操作选项的索引
			minCursor := separatorIndex + 1
			if m.cursor < minCursor {
				m.cursor = totalItems - 1
			}
		case "down", "j":
			m.cursor++
			// 光标最大值为最后一个操作选项的索引
			if m.cursor >= totalItems {
				m.cursor = separatorIndex + 1
			}
		case "enter":
			// 只处理操作选项
			actionIndex := m.cursor - separatorIndex - 1
			// 保存主菜单操作选项的相对索引
			m.mainMenuActionIndex = actionIndex
			switch actionIndex {
			case 0: // 切换环境变量组合
				m.state = SwitchGroupView
				m.cursor = 0
			case 1: // 测速
				m.state = SpeedTestView
				m.cursor = 0
				m.speedTestResults = make(map[string]speedtest.TestResult)
				m.speedTestInProgress = true
				// 返回命令以启动测速
				return m, func() tea.Msg {
					return speedTestStartMsg{}
				}
			case 2: // 编辑组合
				m.state = EditGroupSelectView
				m.cursor = 0
			case 3: // 添加新组合
				m.state = AddGroupView
				m.inputStep = 0
				m.newGroupName = ""
				m.newBaseURL = ""
				m.newAuthToken = ""
				m.activateNew = false
				m.textInput.SetValue("")
				m.textInput.Placeholder = "请输入组合名称"
				m.textInput.Focus() // 确保 textInput 获得焦点
			case 4: // 删除组合
				m.state = DeleteGroupView
				m.cursor = 0
			case 5: // 退出
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// viewMainMenu 渲染主菜单
func (m Model) viewMainMenu() string {
	groups := m.manager.GetGroups()
	separatorIndex := len(groups)

	s := titleStyle.Render("🔧 环境变量管理器") + "\n\n"

	// 显示组合列表（仅供展示，无光标）
	if len(groups) > 0 {
		for _, group := range groups {
			// 激活状态标识
			status := ""
			if group.IsActive {
				status = activeStatusStyle.Render(" ✓") + activeStatusStyle.Render(" (当前)")
			}

			s += menuItemStyle.Render(fmt.Sprintf("  %s%s", group.Name, status)) + "\n"
		}

		// 分隔线
		s += "\n" + subtleStyle.Render("────────────────") + "\n\n"
	}

	// 显示操作选项
	for i, option := range mainMenuOptions {
		itemIndex := separatorIndex + 1 + i
		cursor := " "
		style := menuItemStyle
		if m.cursor == itemIndex {
			cursor = ">"
			style = selectedMenuItemStyle
		}
		s += style.Render(fmt.Sprintf("%s %s", cursor, option)) + "\n"
	}

	s += "\n" + helpStyle.Render("↑/↓: 移动 | Enter: 选择 | q: 退出")

	return s
}
