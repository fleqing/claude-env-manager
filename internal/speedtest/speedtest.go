package speedtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TestResult 存储测速结果
type TestResult struct {
	Success  bool          // 是否成功
	Latency  time.Duration // 延迟时间
	Error    string        // 错误信息
	Duration time.Duration // 请求总耗时
}

// TestGroup 测试单个 API 组合的连接性和延迟
func TestGroup(baseURL, apiKey string) TestResult {
	start := time.Now()

	// 构建请求体
	requestBody := map[string]interface{}{
		"model":      "claude-sonnet-4-5-20250929",
		"max_tokens": 1,
		"messages": []map[string]string{
			{"role": "user", "content": "hi"},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return TestResult{
			Success: false,
			Error:   "构建请求失败",
		}
	}

	// 创建 HTTP 请求
	url := baseURL + "/v1/messages"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return TestResult{
			Success: false,
			Error:   "创建请求失败",
		}
	}

	// 设置请求头
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	// 创建 HTTP 客户端，设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		// 判断错误类型
		errorMsg := "连接失败"
		if err.Error() != "" {
			if contains(err.Error(), "timeout") {
				errorMsg = "超时"
			} else if contains(err.Error(), "connection refused") {
				errorMsg = "连接被拒绝"
			}
		}
		return TestResult{
			Success:  false,
			Error:    errorMsg,
			Duration: duration,
		}
	}
	defer resp.Body.Close()

	// 读取响应体（即使不使用，也要读取以完成请求）
	_, _ = io.ReadAll(resp.Body)

	// 检查 HTTP 状态码
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return TestResult{
			Success:  false,
			Error:    "认证失败",
			Duration: duration,
		}
	}

	if resp.StatusCode >= 400 {
		return TestResult{
			Success:  false,
			Error:    fmt.Sprintf("HTTP %d", resp.StatusCode),
			Duration: duration,
		}
	}

	// 成功
	return TestResult{
		Success:  true,
		Latency:  duration,
		Duration: duration,
	}
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
