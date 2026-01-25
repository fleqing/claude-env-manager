package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// updateSwitchGroup 更新切换组合视图
func (m Model) updateSwitchGroup(msg tea.Msg) (tea.Model, tea.Cmd) {
	groups := m.manager.GetGroups()
	maxCursor := len(groups) // 包含"返回上一级"选项

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor == 0 {
				m.cursor = maxCursor // 循环到最后一项
			} else {
				m.cursor--
			}
		case "down", "j":
			if m.cursor == maxCursor {
				m.cursor = 0 // 循环到第一项
			} else {
				m.cursor++
			}
		case "enter":
			if m.cursor == len(groups) {
				// 选择了"返回上一级"
				m.state = MainMenuView
				groups := m.manager.GetGroups()
				m.cursor = len(groups) + 1 + m.mainMenuActionIndex
			} else {
				// 选择了某个组合，直接切换
				selectedGroupName := groups[m.cursor].Name
				if err := m.manager.ActivateGroup(selectedGroupName); err != nil {
					m.err = err
					m.message = ""
				} else {
					m.message = successStyle.Render(fmt.Sprintf("✓ 已切换到 %s", selectedGroupName)) + "\n\n" +
						warningStyle.Render("💡 请执行: source ~/.zshrc")
					m.err = nil
				}
				m.state = MainMenuView
				groups := m.manager.GetGroups()
				m.cursor = len(groups) + 1 + m.mainMenuActionIndex
			}
		}
	}
	return m, nil
}

// viewSwitchGroup 渲染切换组合视图
func (m Model) viewSwitchGroup() string {
	groups := m.manager.GetGroups()

	s := titleStyle.Render("🔄 切换环境变量组合") + "\n\n"

	if len(groups) == 0 {
		s += subtleStyle.Render("暂无环境变量组合") + "\n"
		s += "\n" + helpStyle.Render("ESC: 返回主菜单")
		return s
	}

	s += "请选择要激活的组合：\n\n"

	for i, group := range groups {
		cursor := " "
		style := menuItemStyle
		if m.cursor == i {
			cursor = ">"
			style = selectedMenuItemStyle
		}

		status := ""
		if group.IsActive {
			status = activeStatusStyle.Render(" (当前激活)")
		}

		s += style.Render(fmt.Sprintf("%s %s%s", cursor, group.Name, status)) + "\n"
	}

	// 添加"返回上一级"选项
	cursor := " "
	style := menuItemStyle
	if m.cursor == len(groups) {
		cursor = ">"
		style = selectedMenuItemStyle
	}
	s += "\n" + style.Render(fmt.Sprintf("%s 返回上一级", cursor)) + "\n"

	s += "\n" + helpStyle.Render("↑/↓: 移动 | Enter: 切换 | ESC: 返回")

	return s
}
