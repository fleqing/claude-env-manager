package ui

import (
	"claude-env-manager/internal/model"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// updateAddGroup 更新添加组合视图
func (m Model) updateAddGroup(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			value := m.textInput.Value()
			if value == "" {
				m.err = fmt.Errorf("输入不能为空")
				return m, nil
			}

			// 根据当前步骤保存输入
			switch m.inputStep {
			case 0: // 输入名称
				m.newGroupName = value
				m.inputStep++
				m.textInput.SetValue("")
				m.textInput.Placeholder = "请输入 ANTHROPIC_BASE_URL"
			case 1: // 输入 BASE_URL
				m.newBaseURL = value
				m.inputStep++
				m.textInput.SetValue("")
				m.textInput.Placeholder = "请输入 ANTHROPIC_AUTH_TOKEN"
			case 2: // 输入 AUTH_TOKEN
				m.newAuthToken = value
				m.inputStep++
				m.textInput.Blur()
				m.cursor = 0 // 默认选择"是"
			case 3: // 确认是否激活
				// 根据光标位置设置是否激活
				m.activateNew = (m.cursor == 0)

				// 创建新组合
				newGroup := model.EnvGroup{
					Name:      m.newGroupName,
					BaseURL:   m.newBaseURL,
					AuthToken: m.newAuthToken,
					IsActive:  m.activateNew,
				}

				if err := m.manager.AddGroup(newGroup); err != nil {
					m.err = err
					m.message = ""
				} else {
					m.message = successStyle.Render(fmt.Sprintf("✓ 已添加组合 %s", m.newGroupName)) + "\n\n" +
						warningStyle.Render("💡 请执行: source ~/.zshrc")
					m.err = nil
				}

				// 重置状态
				m.state = MainMenuView
				groups := m.manager.GetGroups()
				m.cursor = len(groups) + 1 + m.mainMenuActionIndex
				m.inputStep = 0
				m.newGroupName = ""
				m.newBaseURL = ""
				m.newAuthToken = ""
				m.activateNew = false
				m.textInput.SetValue("")
				return m, nil
			}
			return m, nil
		case "up", "k":
			// 在确认步骤时，允许上下键选择
			if m.inputStep == 3 {
				m.cursor = (m.cursor + 1) % 2
			}
		case "down", "j":
			// 在确认步骤时，允许上下键选择
			if m.inputStep == 3 {
				m.cursor = (m.cursor + 1) % 2
			}
		case "esc":
			// ESC 返回主菜单
			m.state = MainMenuView
			groups := m.manager.GetGroups()
			m.cursor = len(groups) + 1 + m.mainMenuActionIndex
			m.inputStep = 0
			m.newGroupName = ""
			m.newBaseURL = ""
			m.newAuthToken = ""
			m.activateNew = false
			m.textInput.SetValue("")
			m.err = nil
			return m, nil
		}
	}

	// 只在前三步更新文本输入
	if m.inputStep < 3 {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

// viewAddGroup 渲染添加组合视图
func (m Model) viewAddGroup() string {
	s := titleStyle.Render("➕ 添加新组合") + "\n\n"

	switch m.inputStep {
	case 0:
		s += "请输入组合名称：\n\n"
		s += m.textInput.View() + "\n\n"
	case 1:
		s += fmt.Sprintf("组合名称: %s\n\n", successStyle.Render(m.newGroupName))
		s += "请输入 ANTHROPIC_BASE_URL：\n\n"
		s += m.textInput.View() + "\n\n"
	case 2:
		s += fmt.Sprintf("组合名称: %s\n", successStyle.Render(m.newGroupName))
		s += fmt.Sprintf("BASE_URL: %s\n\n", successStyle.Render(m.newBaseURL))
		s += "请输入 ANTHROPIC_AUTH_TOKEN：\n\n"
		s += m.textInput.View() + "\n\n"
	case 3:
		s += fmt.Sprintf("组合名称: %s\n", successStyle.Render(m.newGroupName))
		s += fmt.Sprintf("BASE_URL: %s\n", successStyle.Render(m.newBaseURL))
		s += fmt.Sprintf("AUTH_TOKEN: %s\n\n", successStyle.Render(m.newAuthToken))
		s += "是否立即激活此组合？\n\n"

		// 显示选项
		options := []string{"是", "否"}
		for i, option := range options {
			cursor := " "
			style := menuItemStyle
			if m.cursor == i {
				cursor = ">"
				style = selectedMenuItemStyle
			}
			s += style.Render(fmt.Sprintf("%s %s", cursor, option)) + "\n"
		}
		s += "\n"
	}

	if m.err != nil {
		s += errorStyle.Render(fmt.Sprintf("✗ %s", m.err.Error())) + "\n\n"
	}

	if m.inputStep == 3 {
		s += helpStyle.Render("↑/↓: 选择 | Enter: 确认 | ESC: 取消并返回")
	} else {
		s += helpStyle.Render("Enter: 下一步 | ESC: 取消并返回")
	}

	return s
}
