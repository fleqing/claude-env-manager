package ui

import (
	"claude-env-manager/internal/speedtest"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// updateSpeedTest 更新测速视图
func (m Model) updateSpeedTest(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// 返回主菜单
			return m.handleBack(), nil
		}
	case speedTestStartMsg:
		// 开始测速所有组合
		groups := m.manager.GetGroups()
		var cmds []tea.Cmd
		for _, group := range groups {
			// 为每个组合创建一个测速命令
			groupName := group.Name
			baseURL := group.BaseURL
			apiKey := group.AuthToken
			cmds = append(cmds, func() tea.Msg {
				result := speedtest.TestGroup(baseURL, apiKey)
				return speedTestResultMsg{
					groupName: groupName,
					result:    result,
				}
			})
		}
		return m, tea.Batch(cmds...)
	case speedTestResultMsg:
		// 更新测速结果
		m.speedTestResults[msg.groupName] = msg.result
		// 检查是否所有测速都完成
		groups := m.manager.GetGroups()
		if len(m.speedTestResults) >= len(groups) {
			m.speedTestInProgress = false
		}
		return m, nil
	}
	return m, nil
}

// viewSpeedTest 渲染测速视图
func (m Model) viewSpeedTest() string {
	s := titleStyle.Render("⚡ 测速结果") + "\n\n"

	groups := m.manager.GetGroups()

	if len(groups) == 0 {
		s += subtleStyle.Render("没有可测速的组合") + "\n"
	} else {
		for _, group := range groups {
			result, exists := m.speedTestResults[group.Name]

			if !exists {
				// 测试中
				s += menuItemStyle.Render(fmt.Sprintf("  %s [测试中...]", group.Name)) + "\n"
			} else if result.Success {
				// 成功 - 显示延迟
				latencyMs := result.Latency.Milliseconds()
				s += successStyle.Render(fmt.Sprintf("  %s [%dms]", group.Name, latencyMs)) + "\n"
			} else {
				// 失败 - 显示错误
				s += errorStyle.Render(fmt.Sprintf("  %s [✗ %s]", group.Name, result.Error)) + "\n"
			}
		}
	}

	s += "\n" + helpStyle.Render("ESC: 返回主菜单")

	return s
}
