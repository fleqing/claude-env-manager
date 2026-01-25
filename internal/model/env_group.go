package model

// EnvGroup 表示一个环境变量组合
type EnvGroup struct {
	Name      string // 组合名称
	BaseURL   string // ANTHROPIC_BASE_URL 的值
	AuthToken string // ANTHROPIC_AUTH_TOKEN 的值
	IsActive  bool   // 是否激活（未被注释）
	LineStart int    // 在文件中的起始行号
	LineEnd   int    // 在文件中的结束行号
}

// TruncateToken 截断 token 用于显示
func (g *EnvGroup) TruncateToken(maxLen int) string {
	if len(g.AuthToken) <= maxLen {
		return g.AuthToken
	}
	return g.AuthToken[:maxLen] + "..."
}
