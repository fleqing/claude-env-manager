package model

import "testing"

func TestTruncateToken_ShortToken(t *testing.T) {
	g := &EnvGroup{AuthToken: "sk-abc"}
	result := g.TruncateToken(10)
	if result != "sk-abc" {
		t.Errorf("期望 %q，得到 %q", "sk-abc", result)
	}
}

func TestTruncateToken_ExactLength(t *testing.T) {
	g := &EnvGroup{AuthToken: "1234567890"}
	result := g.TruncateToken(10)
	if result != "1234567890" {
		t.Errorf("期望 %q，得到 %q", "1234567890", result)
	}
}

func TestTruncateToken_LongToken(t *testing.T) {
	g := &EnvGroup{AuthToken: "sk-ant-very-long-token"}
	result := g.TruncateToken(8)
	if result != "sk-ant-v..." {
		t.Errorf("期望 %q，得到 %q", "sk-ant-v...", result)
	}
}

func TestTruncateToken_EmptyToken(t *testing.T) {
	g := &EnvGroup{AuthToken: ""}
	result := g.TruncateToken(10)
	if result != "" {
		t.Errorf("期望空字符串，得到 %q", result)
	}
}

func TestTruncateToken_ZeroMaxLen(t *testing.T) {
	g := &EnvGroup{AuthToken: "sk-abc"}
	result := g.TruncateToken(0)
	// maxLen=0 时，len("sk-abc")=6 > 0，截断为空前缀 + "..."
	if result != "..." {
		t.Errorf("期望 %q，得到 %q", "...", result)
	}
}

func TestEnvGroup_Fields(t *testing.T) {
	g := EnvGroup{
		Name:      "测试组合",
		BaseURL:   "https://api.example.com",
		AuthToken: "sk-test-token",
		IsActive:  true,
		LineStart: 5,
		LineEnd:   7,
	}
	if g.Name != "测试组合" {
		t.Errorf("Name 字段不符，期望 %q，得到 %q", "测试组合", g.Name)
	}
	if !g.IsActive {
		t.Errorf("IsActive 应为 true")
	}
	if g.LineEnd-g.LineStart != 2 {
		t.Errorf("LineEnd - LineStart 应为 2，得到 %d", g.LineEnd-g.LineStart)
	}
}

func TestEnvGroup_InactiveByDefault(t *testing.T) {
	g := EnvGroup{
		Name:      "默认组合",
		BaseURL:   "https://api.example.com",
		AuthToken: "sk-token",
	}
	if g.IsActive {
		t.Errorf("未显式设置时 IsActive 应为 false")
	}
}
