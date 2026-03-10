package speedtest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestTestGroup_Success 测试 HTTP 200 成功响应
func TestTestGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Errorf("期望路径 /v1/messages，得到 %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("期望 POST 方法，得到 %s", r.Method)
		}
		if r.Header.Get("x-api-key") == "" {
			t.Errorf("缺少 x-api-key 请求头")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := TestGroup(server.URL, "sk-test-key")

	if !result.Success {
		t.Errorf("HTTP 200 期望成功，得到失败，错误: %s", result.Error)
	}
	if result.Duration <= 0 {
		t.Errorf("Duration 应大于 0")
	}
	if result.Error != "" {
		t.Errorf("成功时 Error 应为空，得到: %s", result.Error)
	}
}

// TestTestGroup_Unauthorized 测试 HTTP 401 认证失败
func TestTestGroup_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	result := TestGroup(server.URL, "invalid-key")

	if result.Success {
		t.Errorf("HTTP 401 应返回失败")
	}
	if result.Error != "认证失败" {
		t.Errorf("HTTP 401 错误信息应为 '认证失败'，得到: %s", result.Error)
	}
}

// TestTestGroup_Forbidden 测试 HTTP 403 禁止访问
func TestTestGroup_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	result := TestGroup(server.URL, "forbidden-key")

	if result.Success {
		t.Errorf("HTTP 403 应返回失败")
	}
	if result.Error != "认证失败" {
		t.Errorf("HTTP 403 错误信息应为 '认证失败'，得到: %s", result.Error)
	}
}

// TestTestGroup_ServerError 测试 HTTP 500 服务端错误
func TestTestGroup_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	result := TestGroup(server.URL, "sk-key")

	if result.Success {
		t.Errorf("HTTP 500 应返回失败")
	}
	if result.Error != "HTTP 500" {
		t.Errorf("期望错误 'HTTP 500'，得到: %s", result.Error)
	}
}

// TestTestGroup_ConnectionRefused 测试连接被拒绝
func TestTestGroup_ConnectionRefused(t *testing.T) {
	result := TestGroup("http://127.0.0.1:1", "sk-key")

	if result.Success {
		t.Errorf("连接失败时应返回失败")
	}
	if result.Error == "" {
		t.Errorf("连接失败时 Error 不应为空")
	}
}

// TestTestGroup_DurationMeasured 测试延迟时间被正确测量
func TestTestGroup_DurationMeasured(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := TestGroup(server.URL, "sk-key")

	if result.Duration < 10*time.Millisecond {
		t.Errorf("Duration 应至少 10ms，得到: %v", result.Duration)
	}
}

// TestTestGroup_LatencySetOnSuccess 测试成功时 Latency 等于 Duration
func TestTestGroup_LatencySetOnSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := TestGroup(server.URL, "sk-key")

	if result.Success && result.Latency != result.Duration {
		t.Errorf("成功时 Latency 应等于 Duration，Latency=%v, Duration=%v", result.Latency, result.Duration)
	}
}

// TestContains_Match 测试 contains 找到子串
func TestContains_Match(t *testing.T) {
	if !contains("connection refused error", "connection refused") {
		t.Errorf("contains 应找到子串 'connection refused'")
	}
}

// TestContains_NoMatch 测试 contains 找不到子串
func TestContains_NoMatch(t *testing.T) {
	if contains("hello world", "timeout") {
		t.Errorf("contains 不应找到不存在的子串 'timeout'")
	}
}

// TestContains_EqualStrings 测试 contains 相等字符串
func TestContains_EqualStrings(t *testing.T) {
	if !contains("timeout", "timeout") {
		t.Errorf("相等字符串 contains 应返回 true")
	}
}

// TestContains_SubstrLonger 测试 substr 比 s 更长时返回 false
func TestContains_SubstrLonger(t *testing.T) {
	if contains("hi", "hello world") {
		t.Errorf("substr 比 s 更长时 contains 应返回 false")
	}
}

// TestContainsHelper_Basic 测试 containsHelper 基本功能
func TestContainsHelper_Basic(t *testing.T) {
	if !containsHelper("abcdef", "cde") {
		t.Errorf("containsHelper 应找到子串 'cde'")
	}
	if containsHelper("abcdef", "xyz") {
		t.Errorf("containsHelper 不应找到 'xyz'")
	}
}
