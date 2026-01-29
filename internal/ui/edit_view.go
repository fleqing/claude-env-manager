package ui

import (
	"claude-env-manager/internal/model"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// 编辑字段选项
var editFieldOptions = []string{
	"名称",
	"BASE_URL",
	"AUTH_TOKEN",
	"返回上一级",
}

// updateEditGroupSelect 更新编辑组合选择视图
func (m Model) updateEditGroupSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				// 选择了某个组合
				m.editGroupName = groups[m.cursor].Name
				m.state = EditFieldSelectView
				m.cursor = 0
			}
		}
	}
	return m, nil
}

// viewEditGroupSelect 渲染编辑组合选择视图
func (m Model) viewEditGroupSelect() string {
	groups := m.manager.GetGroups()

	s := titleStyle.Render("✏️  编辑组合") + "\n\n"

	if len(groups) == 0 {
		s += subtleStyle.Render("暂无环境变量组合") + "\n"
		s += "\n" + helpStyle.Render("ESC: 返回主菜单")
		return s
	}

	s += "请选择要编辑的组合：\n\n"

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

// updateEditFieldSelect 更新编辑字段选择视图
func (m Model) updateEditFieldSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor == 0 {
				m.cursor = len(editFieldOptions) - 1 // 循环到最后一项
			} else {
				m.cursor--
			}
		case "down", "j":
			if m.cursor == len(editFieldOptions)-1 {
				m.cursor = 0 // 循环到第一项
			} else {
				m.cursor++
			}
		case "enter":
			if m.cursor == len(editFieldOptions)-1 {
				// 选择了"返回上一级"
				m.state = EditGroupSelectView
				m.cursor = 0
			} else {
				// 选择了某个字段
				m.editField = editFieldOptions[m.cursor]
				m.state = EditInputView

				// 获取当前组合的信息，设置为默认值
				groups := m.manager.GetGroups()
				for _, g := range groups {
					if g.Name == m.editGroupName {
						switch m.cursor {
						case 0: // 名称
							m.textInput.SetValue(g.Name)
							m.textInput.Placeholder = "请输入新名称"
						case 1: // BASE_URL
							m.textInput.SetValue(g.BaseURL)
							m.textInput.Placeholder = "请输入新的 BASE_URL"
						case 2: // AUTH_TOKEN
							m.textInput.SetValue(g.AuthToken)
							m.textInput.Placeholder = "请输入新的 AUTH_TOKEN"
						}
						break
					}
				}
				m.textInput.Focus()
			}
		}
	}
	return m, nil
}

// viewEditFieldSelect 渲染编辑字段选择视图
func (m Model) viewEditFieldSelect() string {
	s := titleStyle.Render(fmt.Sprintf("✏️  编辑组合: %s", m.editGroupName)) + "\n\n"
	s += "请选择要编辑的字段：\n\n"

	for i, option := range editFieldOptions {
		cursor := " "
		style := menuItemStyle
		if m.cursor == i {
			cursor = ">"
			style = selectedMenuItemStyle
		}

		s += style.Render(fmt.Sprintf("%s %s", cursor, option)) + "\n"
	}

	s += "\n" + helpStyle.Render("↑/↓: 移动 | Enter: 选择 | ESC: 返回")

	return s
}

// updateEditInput 更新编辑输入视图
func (m Model) updateEditInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// 执行更新
			newValue := m.textInput.Value()
			if newValue == "" {
				m.err = fmt.Errorf("输入不能为空")
				return m, nil
			}

			// 获取当前组合
			groups := m.manager.GetGroups()
			var currentGroup model.EnvGroup
			for _, g := range groups {
				if g.Name == m.editGroupName {
					currentGroup = g
					break
				}
			}

			// 根据编辑的字段更新
			switch m.editField {
			case "名称":
				currentGroup.Name = newValue
			case "BASE_URL":
				currentGroup.BaseURL = newValue
			case "AUTH_TOKEN":
				currentGroup.AuthToken = newValue
			}

			// 保存更新
			if err := m.manager.UpdateGroup(m.editGroupName, currentGroup); err != nil {
				m.err = err
				m.message = ""
			} else {
				m.message = successStyle.Render(fmt.Sprintf("✓ 已更新组合 %s 的 %s", m.editGroupName, m.editField)) + "\n\n" +
					warningStyle.Render("💡 请执行: source ~/.zshrc")
				m.err = nil
			}

			m.state = MainMenuView
			groups = m.manager.GetGroups()
			m.cursor = len(groups) + 1 + m.mainMenuActionIndex
			m.textInput.SetValue("")
			return m, nil
		case "esc":
			// 取消编辑
			m.state = EditFieldSelectView
			m.cursor = 0
			m.textInput.SetValue("")
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// viewEditInput 渲染编辑输入视图
func (m Model) viewEditInput() string {
	s := titleStyle.Render(fmt.Sprintf("✏️  编辑 %s: %s", m.editGroupName, m.editField)) + "\n\n"
	s += fmt.Sprintf("请输入新的 %s：\n\n", m.editField)
	s += m.textInput.View() + "\n\n"

	if m.err != nil {
		s += errorStyle.Render(fmt.Sprintf("✗ %s", m.err.Error())) + "\n\n"
	}

	s += helpStyle.Render("Enter: 确认 | ESC: 取消")

	return s
}
