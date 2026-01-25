package ui

import (
	"claude-env-manager/internal/manager"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ViewState 视图状态
type ViewState int

const (
	MainMenuView ViewState = iota
	ListGroupsView
	SwitchGroupView
	SwitchConfirmView
	EditGroupSelectView
	EditFieldSelectView
	EditInputView
	AddGroupView
	AddGroupInputView
	DeleteGroupView
	DeleteConfirmView
)

// InputField 输入字段类型
type InputField int

const (
	InputName InputField = iota
	InputBaseURL
	InputAuthToken
	InputConfirm
)

// Model Bubbletea 主模型
type Model struct {
	manager        *manager.Manager
	state               ViewState
	prevState           ViewState // 用于返回上一级
	cursor              int
	mainMenuActionIndex int // 保存主菜单操作选项的相对索引（0-4）
	message             string
	err                 error

	// 选择相关
	selectedGroup string
	selectedField string

	// 输入相关
	textInput    textinput.Model
	inputField   InputField
	inputStep    int
	newGroupName string
	newBaseURL   string
	newAuthToken string
	activateNew  bool

	// 编辑相关
	editGroupName string
	editField     string

	// 确认相关
	confirmAction string
}

// NewModel 创建新的模型实例
func NewModel(mgr *manager.Manager) Model {
	ti := textinput.New()
	ti.Placeholder = "请输入..."
	ti.Focus()

	// 设置初始光标到第一个操作选项的位置
	groups := mgr.GetGroups()
	initialCursor := len(groups) + 1 // 组合数量 + 分隔线 + 1

	return Model{
		manager:             mgr,
		state:               MainMenuView,
		cursor:              initialCursor,
		mainMenuActionIndex: 0, // 默认选中第一个操作选项（切换环境变量组合）
		textInput:           ti,
	}
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update 更新模型
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == MainMenuView {
				return m, tea.Quit
			}
		case "esc":
			// ESC 键返回上一级或主菜单
			return m.handleBack(), nil
		}
	}

	// 根据当前状态分发更新
	switch m.state {
	case MainMenuView:
		return m.updateMainMenu(msg)
	case ListGroupsView:
		return m.updateListGroups(msg)
	case SwitchGroupView:
		return m.updateSwitchGroup(msg)
	case EditGroupSelectView:
		return m.updateEditGroupSelect(msg)
	case EditFieldSelectView:
		return m.updateEditFieldSelect(msg)
	case EditInputView:
		return m.updateEditInput(msg)
	case AddGroupView:
		return m.updateAddGroup(msg)
	case DeleteGroupView:
		return m.updateDeleteGroup(msg)
	case DeleteConfirmView:
		return m.updateDeleteConfirm(msg)
	}

	return m, nil
}

// View 渲染视图
func (m Model) View() string {
	var s string

	switch m.state {
	case MainMenuView:
		s = m.viewMainMenu()
	case ListGroupsView:
		s = m.viewListGroups()
	case SwitchGroupView:
		s = m.viewSwitchGroup()
	case EditGroupSelectView:
		s = m.viewEditGroupSelect()
	case EditFieldSelectView:
		s = m.viewEditFieldSelect()
	case EditInputView:
		s = m.viewEditInput()
	case AddGroupView:
		s = m.viewAddGroup()
	case DeleteGroupView:
		s = m.viewDeleteGroup()
	case DeleteConfirmView:
		s = m.viewDeleteConfirm()
	}

	// 显示消息或错误
	if m.message != "" {
		s = successStyle.Render(m.message) + "\n\n" + s
	}
	if m.err != nil {
		s = errorStyle.Render(fmt.Sprintf("✗ 错误: %s", m.err.Error())) + "\n\n" + s
	}

	return s
}

// handleBack 处理返回上一级
func (m Model) handleBack() Model {
	switch m.state {
	case ListGroupsView, SwitchGroupView, EditGroupSelectView, AddGroupView, DeleteGroupView:
		m.state = MainMenuView
		// 恢复主菜单的光标位置
		groups := m.manager.GetGroups()
		m.cursor = len(groups) + 1 + m.mainMenuActionIndex
		m.message = ""
		m.err = nil
	case EditFieldSelectView:
		m.state = EditGroupSelectView
		m.cursor = 0
	case EditInputView:
		m.state = EditFieldSelectView
		m.cursor = 0
		m.textInput.SetValue("")
	case DeleteConfirmView:
		m.state = DeleteGroupView
		m.cursor = 0
	default:
		m.state = MainMenuView
		// 恢复主菜单的光标位置
		groups := m.manager.GetGroups()
		m.cursor = len(groups) + 1 + m.mainMenuActionIndex
	}
	return m
}

