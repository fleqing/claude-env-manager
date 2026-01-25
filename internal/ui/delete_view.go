package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// updateDeleteGroup 更新删除组合视图
func (m Model) updateDeleteGroup(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.cursor = 0
			} else {
				// 选择了某个组合
				m.selectedGroup = groups[m.cursor].Name
				m.state = DeleteConfirmView
				m.cursor = 0 // 默认选择"确认删除"
			}
		}
	}
	return m, nil
}

// viewDeleteGroup 渲染删除组合视图
func (m Model) viewDeleteGroup() string {
	groups := m.manager.GetGroups()

	s := titleStyle.Render("🗑️  删除组合") + "\n\n"

	if len(groups) == 0 {
		s += subtleStyle.Render("暂无环境变量组合") + "\n"
		s += "\n" + helpStyle.Render("ESC: 返回主菜单")
		return s
	}

	s += warningStyle.Render("⚠️  警告：删除操作不可恢复！") + "\n\n"
	s += "请选择要删除的组合：\n\n"

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

	s += "\n" + helpStyle.Render("↑/↓: 移动 | Enter: 选择 | ESC: 返回")

	return s
}

// updateDeleteConfirm 更新删除确认视图
func (m Model) updateDeleteConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.cursor = (m.cursor + 1) % 2
		case "down", "j":
			m.cursor = (m.cursor + 1) % 2
		case "enter":
			if m.cursor == 0 {
				// 确认删除
				if err := m.manager.DeleteGroup(m.selectedGroup); err != nil {
					m.err = err
					m.message = ""
				} else {
					m.message = successStyle.Render(fmt.Sprintf("✓ 已删除组合 %s", m.selectedGroup)) + "\n\n" +
						warningStyle.Render("💡 请执行: source ~/.zshrc")
					m.err = nil
				}
			}
			// 无论确认还是取消，都返回主菜单
			m.state = MainMenuView
			groups := m.manager.GetGroups()
			m.cursor = len(groups) + 1 + m.mainMenuActionIndex
		case "esc":
			// 取消删除
			m.state = DeleteGroupView
			m.cursor = 0
		}
	}
	return m, nil
}

// viewDeleteConfirm 渲染删除确认视图
func (m Model) viewDeleteConfirm() string {
	s := titleStyle.Render("🗑️  确认删除") + "\n\n"
	s += warningStyle.Render("⚠️  警告：此操作不可恢复！") + "\n\n"
	s += fmt.Sprintf("确认删除组合 %s 吗？\n\n", errorStyle.Render(m.selectedGroup))

	// 显示选项
	options := []string{"确认删除", "取消"}
	for i, option := range options {
		cursor := " "
		style := menuItemStyle
		if m.cursor == i {
			cursor = ">"
			style = selectedMenuItemStyle
		}
		s += style.Render(fmt.Sprintf("%s %s", cursor, option)) + "\n"
	}

	s += "\n" + helpStyle.Render("↑/↓: 选择 | Enter: 确认 | ESC: 返回")

	return s
}
