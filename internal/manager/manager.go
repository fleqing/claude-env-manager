package manager

import (
	"claude-env-manager/internal/config"
	"claude-env-manager/internal/model"
	"claude-env-manager/internal/parser"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Manager 环境变量管理器
type Manager struct {
	config *config.Config
	groups []model.EnvGroup
}

// NewManager 创建管理器实例
func NewManager(cfg *config.Config) (*Manager, error) {
	m := &Manager{
		config: cfg,
		groups: []model.EnvGroup{},
	}

	// 加载现有配置
	if err := m.Load(); err != nil {
		return nil, err
	}

	return m, nil
}

// Load 从 .zshrc 文件加载环境变量组合
func (m *Manager) Load() error {
	result, err := parser.ParseZshrc(m.config.ZshrcPath)
	if err != nil {
		return fmt.Errorf("解析文件失败: %w", err)
	}

	m.groups = result.Groups
	return nil
}

// GetGroups 获取所有环境变量组合
func (m *Manager) GetGroups() []model.EnvGroup {
	return m.groups
}

// ActivateGroup 激活指定组合
func (m *Manager) ActivateGroup(name string) error {
	found := false
	for i := range m.groups {
		if m.groups[i].Name == name {
			m.groups[i].IsActive = true
			found = true
		} else {
			m.groups[i].IsActive = false
		}
	}

	if !found {
		return fmt.Errorf("组合 %s 不存在", name)
	}

	return m.Save()
}

// AddGroup 添加新组合
func (m *Manager) AddGroup(group model.EnvGroup) error {
	// 检查名称是否已存在
	for _, g := range m.groups {
		if g.Name == group.Name {
			return fmt.Errorf("组合名称 %s 已存在", group.Name)
		}
	}

	// 如果新组合要激活，则停用其他所有组合
	if group.IsActive {
		for i := range m.groups {
			m.groups[i].IsActive = false
		}
	}

	group.LineStart = -1
	group.LineEnd = -1
	m.groups = append(m.groups, group)

	return m.Save()
}

// UpdateGroup 更新组合信息
func (m *Manager) UpdateGroup(oldName string, newGroup model.EnvGroup) error {
	found := false
	for i := range m.groups {
		if m.groups[i].Name == oldName {
			// 如果修改名称，检查新名称是否与其他组合冲突
			if newGroup.Name != oldName {
				for j, g := range m.groups {
					if j != i && g.Name == newGroup.Name {
						return fmt.Errorf("组合名称 %s 已存在", newGroup.Name)
					}
				}
			}

			m.groups[i].Name = newGroup.Name
			m.groups[i].BaseURL = newGroup.BaseURL
			m.groups[i].AuthToken = newGroup.AuthToken
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("组合 %s 不存在", oldName)
	}

	return m.Save()
}

// DeleteGroup 删除指定组合
func (m *Manager) DeleteGroup(name string) error {
	newGroups := []model.EnvGroup{}
	found := false

	for _, g := range m.groups {
		if g.Name != name {
			newGroups = append(newGroups, g)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("组合 %s 不存在", name)
	}

	m.groups = newGroups
	return m.Save()
}

// Save 保存修改到 .zshrc 文件
func (m *Manager) Save() error {
	// 创建备份
	if err := m.createBackup(); err != nil {
		return fmt.Errorf("创建备份失败: %w", err)
	}

	// 读取原文件内容
	result, err := parser.ParseZshrc(m.config.ZshrcPath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	// 构建新文件内容
	var newLines []string

	// 添加其他非相关行，并去除末尾的连续空行
	otherLines := result.OtherLines
	// 从末尾开始，移除所有连续的空行
	for len(otherLines) > 0 && strings.TrimSpace(otherLines[len(otherLines)-1]) == "" {
		otherLines = otherLines[:len(otherLines)-1]
	}
	newLines = append(newLines, otherLines...)

	// 添加环境变量组合
	for _, group := range m.groups {
		newLines = append(newLines, "")
		newLines = append(newLines, fmt.Sprintf("# %s", group.Name))

		if group.IsActive {
			newLines = append(newLines, fmt.Sprintf("export ANTHROPIC_BASE_URL=%s", group.BaseURL))
			newLines = append(newLines, fmt.Sprintf("export ANTHROPIC_AUTH_TOKEN=%s", group.AuthToken))
		} else {
			newLines = append(newLines, fmt.Sprintf("#export ANTHROPIC_BASE_URL=%s", group.BaseURL))
			newLines = append(newLines, fmt.Sprintf("#export ANTHROPIC_AUTH_TOKEN=%s", group.AuthToken))
		}
	}

	// 写入文件
	content := strings.Join(newLines, "\n")
	if err := os.WriteFile(m.config.ZshrcPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	// 重新加载
	return m.Load()
}

// createBackup 创建备份文件
func (m *Manager) createBackup() error {
	// 检查源文件是否存在
	if _, err := os.Stat(m.config.ZshrcPath); os.IsNotExist(err) {
		return nil // 文件不存在，无需备份
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(m.config.BackupDir, fmt.Sprintf("zshrc_backup_%s", timestamp))

	// 复制文件
	input, err := os.ReadFile(m.config.ZshrcPath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(backupPath, input, 0644); err != nil {
		return err
	}

	// 清理旧备份
	return m.cleanOldBackups()
}

// cleanOldBackups 清理旧备份，只保留最近的 N 个
func (m *Manager) cleanOldBackups() error {
	files, err := filepath.Glob(filepath.Join(m.config.BackupDir, "zshrc_backup_*"))
	if err != nil {
		return err
	}

	if len(files) <= m.config.MaxBackups {
		return nil
	}

	// 按修改时间排序
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := os.Stat(files[i])
		infoJ, _ := os.Stat(files[j])
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// 删除最旧的备份
	for i := 0; i < len(files)-m.config.MaxBackups; i++ {
		if err := os.Remove(files[i]); err != nil {
			return err
		}
	}

	return nil
}

