package manager

import (
	"github.com/fleqing/claude-env-manager/internal/config"
	"github.com/fleqing/claude-env-manager/internal/model"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestEnv 创建临时 zshrc 和备份目录，返回 Manager 和清理函数
func setupTestEnv(t *testing.T, zshrcContent string) (*Manager, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "manager_test_*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}

	zshrcPath := filepath.Join(tmpDir, ".zshrc")
	if err := os.WriteFile(zshrcPath, []byte(zshrcContent), 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("创建临时 zshrc 失败: %v", err)
	}

	backupDir := filepath.Join(tmpDir, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("创建备份目录失败: %v", err)
	}

	cfg := &config.Config{
		ZshrcPath:  zshrcPath,
		BackupDir:  backupDir,
		MaxBackups: 3,
	}

	// 直接构造，因为 Manager 字段未导出，且 NewManager 依赖真实 home 目录
	mgr := &Manager{
		config: cfg,
		groups: []model.EnvGroup{},
	}

	if err := mgr.Load(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Load() 失败: %v", err)
	}

	return mgr, func() { os.RemoveAll(tmpDir) }
}

func TestLoad_EmptyFile(t *testing.T) {
	mgr, cleanup := setupTestEnv(t, "")
	defer cleanup()

	groups := mgr.GetGroups()
	if len(groups) != 0 {
		t.Errorf("空文件应返回 0 个组合，得到 %d", len(groups))
	}
}

func TestLoad_WithGroups(t *testing.T) {
	content := "# GroupA\nexport ANTHROPIC_BASE_URL=https://a.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-aaa\n"
	mgr, cleanup := setupTestEnv(t, content)
	defer cleanup()

	groups := mgr.GetGroups()
	if len(groups) != 1 {
		t.Fatalf("期望 1 个组合，得到 %d", len(groups))
	}
	if groups[0].Name != "GroupA" {
		t.Errorf("组合名称应为 GroupA，得到 %s", groups[0].Name)
	}
}

func TestAddGroup_Success(t *testing.T) {
	mgr, cleanup := setupTestEnv(t, "")
	defer cleanup()

	newGroup := model.EnvGroup{
		Name:      "NewGroup",
		BaseURL:   "https://new.example.com",
		AuthToken: "sk-new",
		IsActive:  false,
	}

	if err := mgr.AddGroup(newGroup); err != nil {
		t.Fatalf("AddGroup 不应返回错误: %v", err)
	}

	groups := mgr.GetGroups()
	if len(groups) != 1 {
		t.Fatalf("添加后应有 1 个组合，得到 %d", len(groups))
	}
	if groups[0].Name != "NewGroup" {
		t.Errorf("组合名称不匹配，期望 NewGroup，得到 %s", groups[0].Name)
	}
}

func TestAddGroup_DuplicateName(t *testing.T) {
	content := "# ExistingGroup\nexport ANTHROPIC_BASE_URL=https://api.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-existing\n"
	mgr, cleanup := setupTestEnv(t, content)
	defer cleanup()

	err := mgr.AddGroup(model.EnvGroup{
		Name:      "ExistingGroup",
		BaseURL:   "https://other.example.com",
		AuthToken: "sk-other",
	})
	if err == nil {
		t.Errorf("添加重名组合应返回错误")
	}
}

func TestAddGroup_ActiveDeactivatesOthers(t *testing.T) {
	content := "# GroupA\nexport ANTHROPIC_BASE_URL=https://a.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-aaa\n"
	mgr, cleanup := setupTestEnv(t, content)
	defer cleanup()

	if err := mgr.AddGroup(model.EnvGroup{
		Name:      "GroupB",
		BaseURL:   "https://b.example.com",
		AuthToken: "sk-bbb",
		IsActive:  true,
	}); err != nil {
		t.Fatalf("AddGroup 失败: %v", err)
	}

	for _, g := range mgr.GetGroups() {
		if g.Name == "GroupA" && g.IsActive {
			t.Errorf("添加激活的 GroupB 后，GroupA 应被停用")
		}
		if g.Name == "GroupB" && !g.IsActive {
			t.Errorf("GroupB 应为激活状态")
		}
	}
}

func TestDeleteGroup_Success(t *testing.T) {
	content := "# ToDelete\nexport ANTHROPIC_BASE_URL=https://del.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-del\n"
	mgr, cleanup := setupTestEnv(t, content)
	defer cleanup()

	if err := mgr.DeleteGroup("ToDelete"); err != nil {
		t.Fatalf("DeleteGroup 不应返回错误: %v", err)
	}
	if len(mgr.GetGroups()) != 0 {
		t.Errorf("删除后应有 0 个组合")
	}
}

func TestDeleteGroup_NotFound(t *testing.T) {
	mgr, cleanup := setupTestEnv(t, "")
	defer cleanup()

	if err := mgr.DeleteGroup("NonExistent"); err == nil {
		t.Errorf("删除不存在的组合应返回错误")
	}
}

func TestActivateGroup_Success(t *testing.T) {
	content := "# GroupA\n#export ANTHROPIC_BASE_URL=https://a.example.com\n#export ANTHROPIC_AUTH_TOKEN=sk-aaa\n\n# GroupB\n#export ANTHROPIC_BASE_URL=https://b.example.com\n#export ANTHROPIC_AUTH_TOKEN=sk-bbb\n"
	mgr, cleanup := setupTestEnv(t, content)
	defer cleanup()

	if err := mgr.ActivateGroup("GroupA"); err != nil {
		t.Fatalf("ActivateGroup 不应返回错误: %v", err)
	}

	for _, g := range mgr.GetGroups() {
		if g.Name == "GroupA" && !g.IsActive {
			t.Errorf("GroupA 应被激活")
		}
		if g.Name == "GroupB" && g.IsActive {
			t.Errorf("GroupB 应被停用")
		}
	}
}

func TestActivateGroup_NotFound(t *testing.T) {
	mgr, cleanup := setupTestEnv(t, "")
	defer cleanup()

	if err := mgr.ActivateGroup("NonExistent"); err == nil {
		t.Errorf("激活不存在的组合应返回错误")
	}
}

func TestUpdateGroup_Success(t *testing.T) {
	content := "# OldName\nexport ANTHROPIC_BASE_URL=https://old.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-old\n"
	mgr, cleanup := setupTestEnv(t, content)
	defer cleanup()

	updated := model.EnvGroup{
		Name:      "NewName",
		BaseURL:   "https://new.example.com",
		AuthToken: "sk-new",
	}

	if err := mgr.UpdateGroup("OldName", updated); err != nil {
		t.Fatalf("UpdateGroup 不应返回错误: %v", err)
	}

	groups := mgr.GetGroups()
	if len(groups) != 1 {
		t.Fatalf("更新后应仍有 1 个组合")
	}
	if groups[0].Name != "NewName" {
		t.Errorf("更新后名称应为 NewName，得到 %s", groups[0].Name)
	}
	if groups[0].BaseURL != "https://new.example.com" {
		t.Errorf("BaseURL 未更新，得到 %s", groups[0].BaseURL)
	}
}

func TestUpdateGroup_NotFound(t *testing.T) {
	mgr, cleanup := setupTestEnv(t, "")
	defer cleanup()

	err := mgr.UpdateGroup("Ghost", model.EnvGroup{Name: "Ghost", BaseURL: "x", AuthToken: "y"})
	if err == nil {
		t.Errorf("更新不存在的组合应返回错误")
	}
}

func TestUpdateGroup_DuplicateName(t *testing.T) {
	content := "# GroupA\nexport ANTHROPIC_BASE_URL=https://a.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-aaa\n\n# GroupB\nexport ANTHROPIC_BASE_URL=https://b.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-bbb\n"
	mgr, cleanup := setupTestEnv(t, content)
	defer cleanup()

	err := mgr.UpdateGroup("GroupA", model.EnvGroup{
		Name:      "GroupB",
		BaseURL:   "https://a.example.com",
		AuthToken: "sk-aaa",
	})
	if err == nil {
		t.Errorf("改名为已存在的 GroupB 应返回错误")
	}
}

func TestSave_PreservesOtherLines(t *testing.T) {
	content := "export PATH=/usr/local/bin:$PATH\nalias ll='ls -la'\n\n# MyGroup\nexport ANTHROPIC_BASE_URL=https://api.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-token"
	mgr, cleanup := setupTestEnv(t, content)
	defer cleanup()

	if err := mgr.ActivateGroup("MyGroup"); err != nil {
		t.Fatalf("ActivateGroup 失败: %v", err)
	}

	savedContent, err := os.ReadFile(mgr.config.ZshrcPath)
	if err != nil {
		t.Fatalf("读取保存文件失败: %v", err)
	}

	saved := string(savedContent)
	if !strings.Contains(saved, "export PATH=/usr/local/bin:$PATH") {
		t.Errorf("保存后 PATH 设置应被保留")
	}
	if !strings.Contains(saved, "alias ll='ls -la'") {
		t.Errorf("保存后 alias 设置应被保留")
	}
}

func TestGetGroups_ReturnsCorrectCount(t *testing.T) {
	content := "# Group1\nexport ANTHROPIC_BASE_URL=https://api.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-token\n"
	mgr, cleanup := setupTestEnv(t, content)
	defer cleanup()

	groups := mgr.GetGroups()
	if len(groups) != 1 {
		t.Errorf("有 1 个组合时 GetGroups 应返回长度为 1 的切片，得到 %d", len(groups))
	}
}
