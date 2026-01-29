package parser

import (
	"bufio"
	"claude-env-manager/internal/model"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	// 匹配注释行（组合名称）
	commentPattern = regexp.MustCompile(`^#\s*(.+?)\s*$`)
	// 匹配 BASE_URL 行
	baseURLPattern = regexp.MustCompile(`^(#)?export\s+ANTHROPIC_BASE_URL=(.+)$`)
	// 匹配 AUTH_TOKEN 行
	authTokenPattern = regexp.MustCompile(`^(#)?export\s+ANTHROPIC_AUTH_TOKEN=(.+)$`)
)

// ParseResult 解析结果
type ParseResult struct {
	Groups    []model.EnvGroup // 解析出的环境变量组合
	OtherLines []string         // 其他非相关行（用于保留原文件内容）
	LineMap    map[int]bool     // 记录哪些行属于环境变量组合
}

// ParseZshrc 解析 .zshrc 文件
func ParseZshrc(path string) (*ParseResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	var groups []model.EnvGroup
	lineMap := make(map[int]bool)
	var allLines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		allLines = append(allLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 解析环境变量组合
	for i := 0; i < len(allLines); i++ {
		line := strings.TrimSpace(allLines[i])

		// 检查是否是组合名称注释
		if matches := commentPattern.FindStringSubmatch(line); matches != nil {
			// 检查接下来的两行是否是 BASE_URL 和 AUTH_TOKEN
			if i+2 < len(allLines) {
				baseURLLine := strings.TrimSpace(allLines[i+1])
				authTokenLine := strings.TrimSpace(allLines[i+2])

				baseURLMatches := baseURLPattern.FindStringSubmatch(baseURLLine)
				authTokenMatches := authTokenPattern.FindStringSubmatch(authTokenLine)

				if baseURLMatches != nil && authTokenMatches != nil {
					// 找到一个完整的组合
					name := matches[1]
					isCommentedBaseURL := baseURLMatches[1] == "#"
					baseURL := strings.TrimSpace(baseURLMatches[2])
					isCommentedAuthToken := authTokenMatches[1] == "#"
					authToken := strings.TrimSpace(authTokenMatches[2])

					// 只有两行都未被注释才算激活
					isActive := !isCommentedBaseURL && !isCommentedAuthToken

					group := model.EnvGroup{
						Name:      name,
						BaseURL:   baseURL,
						AuthToken: authToken,
						IsActive:  isActive,
						LineStart: i,
						LineEnd:   i + 2,
					}

					groups = append(groups, group)

					// 标记这些行属于环境变量组合
					lineMap[i] = true
					lineMap[i+1] = true
					lineMap[i+2] = true

					// 跳过已处理的行
					i += 2
				}
			}
		}
	}

	// 收集其他非相关行
	var otherLines []string
	for i, line := range allLines {
		if !lineMap[i] {
			otherLines = append(otherLines, line)
		}
	}

	return &ParseResult{
		Groups:     groups,
		OtherLines: otherLines,
		LineMap:    lineMap,
	}, nil
}
