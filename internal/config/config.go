package config

import (
	"os"
	"path/filepath"
)

// Config 存储应用配置
type Config struct {
	ZshrcPath string // .zshrc 文件路径
	BackupDir string // 备份目录路径
	MaxBackups int   // 最大备份数量
}

// NewConfig 创建默认配置
func NewConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	backupDir := filepath.Join(homeDir, ".claude-env-manager", "backups")

	// 确保备份目录存在
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, err
	}

	return &Config{
		ZshrcPath:  filepath.Join(homeDir, ".zshrc"),
		BackupDir:  backupDir,
		MaxBackups: 10,
	}, nil
}
