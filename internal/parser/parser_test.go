package parser

import (
	"os"
	"testing"
)

// writeTemp 在系统临时目录创建内容文件，返回路径，并注册 t.Cleanup 自动删除
func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "zshrc_test_*")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		t.Fatalf("写入临时文件失败: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

// TestParseZshrc_SingleActiveGroup 测试单个激活的组合
func TestParseZshrc_SingleActiveGroup(t *testing.T) {
	content := "# MyGroup\nexport ANTHROPIC_BASE_URL=https://api.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-test-token\n"
	path := writeTemp(t, content)

	result, err := ParseZshrc(path)
	if err != nil {
		t.Fatalf("ParseZshrc 不应返回错误: %v", err)
	}
	if len(result.Groups) != 1 {
		t.Fatalf("期望 1 个组合，得到 %d", len(result.Groups))
	}
	g := result.Groups[0]
	if g.Name != "MyGroup" {
		t.Errorf("Name 期望 %q，得到 %q", "MyGroup", g.Name)
	}
	if g.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL 期望 %q，得到 %q", "https://api.example.com", g.BaseURL)
	}
	if g.AuthToken != "sk-test-token" {
		t.Errorf("AuthToken 期望 %q，得到 %q", "sk-test-token", g.AuthToken)
	}
	if !g.IsActive {
		t.Errorf("组合应为激活状态")
	}
}

// TestParseZshrc_SingleInactiveGroup 测试单个注释（停用）的组合
func TestParseZshrc_SingleInactiveGroup(t *testing.T) {
	content := "# MyGroup\n#export ANTHROPIC_BASE_URL=https://api.example.com\n#export ANTHROPIC_AUTH_TOKEN=sk-test-token\n"
	path := writeTemp(t, content)

	result, err := ParseZshrc(path)
	if err != nil {
		t.Fatalf("ParseZshrc 不应返回错误: %v", err)
	}
	if len(result.Groups) != 1 {
		t.Fatalf("期望 1 个组合，得到 %d", len(result.Groups))
	}
	if result.Groups[0].IsActive {
		t.Errorf("组合应为停用状态")
	}
}

// TestParseZshrc_MixedCommentState 测试仅 BASE_URL 被注释时应为停用
func TestParseZshrc_MixedCommentState(t *testing.T) {
	content := "# HalfGroup\n#export ANTHROPIC_BASE_URL=https://api.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-token\n"
	path := writeTemp(t, content)

	result, err := ParseZshrc(path)
	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if len(result.Groups) != 1 {
		t.Fatalf("期望识别 1 个组合，得到 %d", len(result.Groups))
	}
	if result.Groups[0].IsActive {
		t.Errorf("BASE_URL 被注释时，组合应为停用状态")
	}
}

// TestParseZshrc_MultipleGroups 测试多个组合
func TestParseZshrc_MultipleGroups(t *testing.T) {
	content := "# GroupA\nexport ANTHROPIC_BASE_URL=https://a.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-aaa\n\n# GroupB\n#export ANTHROPIC_BASE_URL=https://b.example.com\n#export ANTHROPIC_AUTH_TOKEN=sk-bbb\n"
	path := writeTemp(t, content)

	result, err := ParseZshrc(path)
	if err != nil {
		t.Fatalf("ParseZshrc 不应返回错误: %v", err)
	}
	if len(result.Groups) != 2 {
		t.Fatalf("期望 2 个组合，得到 %d", len(result.Groups))
	}
	if !result.Groups[0].IsActive {
		t.Errorf("GroupA 应为激活状态")
	}
	if result.Groups[1].IsActive {
		t.Errorf("GroupB 应为停用状态")
	}
}

// TestParseZshrc_PreservesOtherLines 测试其他无关行被保留在 OtherLines 中
func TestParseZshrc_PreservesOtherLines(t *testing.T) {
	content := "export PATH=/usr/local/bin:$PATH\nalias ll='ls -la'\n\n# MyGroup\nexport ANTHROPIC_BASE_URL=https://api.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-token\n"
	path := writeTemp(t, content)

	result, err := ParseZshrc(path)
	if err != nil {
		t.Fatalf("ParseZshrc 不应返回错误: %v", err)
	}
	if len(result.Groups) != 1 {
		t.Fatalf("期望 1 个组合，得到 %d", len(result.Groups))
	}

	foundPath := false
	for _, l := range result.OtherLines {
		if l == "export PATH=/usr/local/bin:$PATH" {
			foundPath = true
		}
	}
	if !foundPath {
		t.Errorf("OtherLines 应包含 PATH 设置行")
	}
}

// TestParseZshrc_EmptyFile 测试空文件
func TestParseZshrc_EmptyFile(t *testing.T) {
	path := writeTemp(t, "")

	result, err := ParseZshrc(path)
	if err != nil {
		t.Fatalf("空文件不应返回错误: %v", err)
	}
	if len(result.Groups) != 0 {
		t.Errorf("空文件应返回 0 个组合，得到 %d", len(result.Groups))
	}
}

// TestParseZshrc_PartialGroup_NoMatch 测试不完整的格式不识别为组合
func TestParseZshrc_PartialGroup_NoMatch(t *testing.T) {
	content := "# JustAComment\nexport PATH=/usr/bin\n"
	path := writeTemp(t, content)

	result, err := ParseZshrc(path)
	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if len(result.Groups) != 0 {
		t.Errorf("不完整的格式不应被识别为组合，得到 %d 个", len(result.Groups))
	}
}

// TestParseZshrc_FileNotExist 测试文件不存在时应返回错误
func TestParseZshrc_FileNotExist(t *testing.T) {
	_, err := ParseZshrc("/tmp/nonexistent_zshrc_test_abc123xyz.zshrc")
	if err == nil {
		t.Errorf("文件不存在时应返回错误")
	}
}

// TestParseZshrc_LineMap 测试 LineMap 正确标记了组合所在行
func TestParseZshrc_LineMap(t *testing.T) {
	content := "# Group1\nexport ANTHROPIC_BASE_URL=https://api.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-token\n"
	path := writeTemp(t, content)

	result, err := ParseZshrc(path)
	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	// 行 0（# Group1）、行 1（BASE_URL）、行 2（AUTH_TOKEN）应被标记
	if !result.LineMap[0] {
		t.Errorf("第 0 行应在 LineMap 中")
	}
	if !result.LineMap[1] {
		t.Errorf("第 1 行应在 LineMap 中")
	}
	if !result.LineMap[2] {
		t.Errorf("第 2 行应在 LineMap 中")
	}
}

// TestParseZshrc_LineStartEnd 测试 LineStart 和 LineEnd 的值正确
func TestParseZshrc_LineStartEnd(t *testing.T) {
	content := "# Group1\nexport ANTHROPIC_BASE_URL=https://api.example.com\nexport ANTHROPIC_AUTH_TOKEN=sk-token\n"
	path := writeTemp(t, content)

	result, err := ParseZshrc(path)
	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if len(result.Groups) != 1 {
		t.Fatalf("期望 1 个组合")
	}
	g := result.Groups[0]
	if g.LineStart != 0 {
		t.Errorf("LineStart 应为 0，得到 %d", g.LineStart)
	}
	if g.LineEnd != 2 {
		t.Errorf("LineEnd 应为 2，得到 %d", g.LineEnd)
	}
}
