package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewConfig_ReturnsValidConfig(t *testing.T) {
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() 不应返回错误，得到: %v", err)
	}
	if cfg == nil {
		t.Fatal("NewConfig() 不应返回 nil")
	}
}

func TestNewConfig_ZshrcPathContainsHome(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("无法获取 home 目录，跳过测试")
	}

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() 失败: %v", err)
	}

	if !strings.HasPrefix(cfg.ZshrcPath, homeDir) {
		t.Errorf("ZshrcPath %q 应以 home 目录 %q 开头", cfg.ZshrcPath, homeDir)
	}
	if filepath.Base(cfg.ZshrcPath) != ".zshrc" {
		t.Errorf("ZshrcPath 文件名应为 .zshrc，得到: %s", filepath.Base(cfg.ZshrcPath))
	}
}

func TestNewConfig_BackupDirExists(t *testing.T) {
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() 失败: %v", err)
	}

	info, err := os.Stat(cfg.BackupDir)
	if err != nil {
		t.Fatalf("备份目录 %q 应存在，但 Stat 失败: %v", cfg.BackupDir, err)
	}
	if !info.IsDir() {
		t.Errorf("备份目录 %q 应为目录", cfg.BackupDir)
	}
}

func TestNewConfig_BackupDirUnderHome(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("无法获取 home 目录，跳过测试")
	}

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() 失败: %v", err)
	}

	if !strings.HasPrefix(cfg.BackupDir, homeDir) {
		t.Errorf("BackupDir %q 应在 home 目录 %q 下", cfg.BackupDir, homeDir)
	}
}

func TestNewConfig_MaxBackupsDefault(t *testing.T) {
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() 失败: %v", err)
	}

	if cfg.MaxBackups != 10 {
		t.Errorf("MaxBackups 默认值应为 10，得到: %d", cfg.MaxBackups)
	}
}
