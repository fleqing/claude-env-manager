package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// updateListGroups 更新查看组合视图
func (m Model) updateListGroups(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "q", "esc":
			m.state = MainMenuView
			groups := m.manager.GetGroups()
			m.cursor = len(groups) + 1 + m.mainMenuActionIndex
		}
	}
	return m, nil
}

// viewListGroups 渲染查看组合视图
func (m Model) viewListGroups() string {
	groups := m.manager.GetGroups()

	s := titleStyle.Render("📋 所有环境变量组合") + "\n\n"

	if len(groups) == 0 {
		s += subtleStyle.Render("暂无环境变量组合") + "\n"
	} else {
		// 表格头部
		header := fmt.Sprintf("%-10s %-20s %-40s %-30s",
			"状态", "名称", "BASE_URL", "AUTH_TOKEN")
		s += tableHeaderStyle.Render(header) + "\n"

		// 表格内容
		for _, group := range groups {
			status := "  "
			statusStyle := inactiveStatusStyle
			if group.IsActive {
				status = "✓ "
				statusStyle = activeStatusStyle
			}

			// 截断过长的 token
			token := group.TruncateToken(25)

			row := fmt.Sprintf("%-10s %-20s %-40s %-30s",
				statusStyle.Render(status),
				group.Name,
				group.BaseURL,
				token)
			s += tableCellStyle.Render(row) + "\n"
		}
	}

	s += "\n" + helpStyle.Render("按任意键返回主菜单")

	return s
}
