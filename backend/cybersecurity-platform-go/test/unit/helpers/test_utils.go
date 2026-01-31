// test/unit/helpers/test_utils.go
package helpers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 创建测试HTTP请求
func CreateTestRequest(method, url string, body interface{}) *http.Request {
	var reqBody io.Reader
	
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	}
	
	return httptest.NewRequest(method, url, reqBody)
}

// 执行HTTP请求并返回响应
func ExecuteRequest(handler http.HandlerFunc, req *http.Request) *http.Response {
	w := httptest.NewRecorder()
	
	// 设置必要的头部
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	handler(w, req)
	return w.Result()
}

// 解析JSON响应
func ParseJSONResponse(t *testing.T, body []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("解析JSON响应失败: %v", err)
	}
	return result
}

// 创建测试文章内容文件
func CreateTestArticleFile(t *testing.T, articleID int, content string) string {
	// 创建临时目录
	tempDir := t.TempDir()
	articlesDir := filepath.Join(tempDir, "forum", "articles")
	
	// 确保目录存在
	err := os.MkdirAll(articlesDir, 0755)
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	
	// 创建测试文件
	filePath := filepath.Join(articlesDir, fmt.Sprintf("%d.txt", articleID))
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	
	return filePath
}

// 验证HTTP响应状态
func AssertStatus(t *testing.T, response *http.Response, expected int) {
	if response.StatusCode != expected {
		body, _ := io.ReadAll(response.Body)
		t.Errorf("期望状态码 %d, 得到 %d。响应: %s", 
			expected, response.StatusCode, string(body))
	}
}

// 验证响应包含特定字段
func AssertResponseField(t *testing.T, response map[string]interface{}, field string) {
	if _, exists := response[field]; !exists {
		t.Errorf("响应中缺少字段: %s", field)
	}
}

// 模拟论坛文章数据
func MockForumArticle() map[string]interface{} {
	return map[string]interface{}{
		"id":         1,
		"title":      "测试文章标题",
		"stuId":      "2023001",
		"stuName":    "测试学生",
		"stuHead":    "default.jpg",
		"cateId":     1,
		"isTop":      false,
		"isEss":      false,
		"viewCount":  100,
		"createTime": "2023-01-01 10:00:00",
		"updateTime": "2023-01-01 10:00:00",
		"content":    "测试文章内容",
	}
}

// 模拟论坛评论数据
func MockForumComment() map[string]interface{} {
	return map[string]interface{}{
		"id":          1,
		"articleId":   1,
		"content":     "测试评论内容",
		"authorId":    "2023001",
		"authorName":  "测试用户",
		"authorHead":  "default.jpg",
		"parentId":    nil,
		"status":      1,
		"createTime":  "2023-01-01 11:00:00",
		"likeCount":   5,
	}
}

// 模拟热门文章数据
func MockHotArticles() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id":         1,
			"title":      "热门文章1",
			"viewCount":  1000,
			"createTime": "2023-01-01 10:00:00",
		},
		{
			"id":         2,
			"title":      "热门文章2",
			"viewCount":  800,
			"createTime": "2023-01-02 10:00:00",
		},
	}
}

// 模拟分类数据
func MockCategories() []map[string]interface{} {
	return []map[string]interface{}{
		{"id": 0, "name": "全部"},
		{"id": 1, "name": "技术讨论"},
		{"id": 2, "name": "学习交流"},
		{"id": 3, "name": "资源共享"},
	}
}

// 模拟标签数据
func MockTags() []map[string]interface{} {
	return []map[string]interface{}{
		{"id": 1, "name": "Go语言", "hot": 50},
		{"id": 2, "name": "网络安全", "hot": 45},
		{"id": 3, "name": "算法", "hot": 40},
		{"id": 4, "name": "数据库", "hot": 35},
		{"id": 5, "name": "前端", "hot": 30},
	}
}

// 验证基础响应格式
func ValidateBaseResponse(t *testing.T, response map[string]interface{}) {
	AssertResponseField(t, response, "code")
	AssertResponseField(t, response, "message")
	AssertResponseField(t, response, "data")
	
	// 验证code是整数
	code, ok := response["code"].(float64) // JSON unmarshal将数字转为float64
	if !ok {
		t.Error("响应中的code字段类型不正确")
	}
	
	// 验证code是有效的（20000表示成功）
	if code != 20000 {
		// 对于错误响应，验证错误信息存在
		message, exists := response["message"].(string)
		if exists && message == "" {
			t.Error("错误响应中缺少错误消息")
		}
	}
}

// 创建测试查询参数
func BuildQueryParams(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	
	var queryParts []string
	for key, value := range params {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, value))
	}
	
	return "?" + strings.Join(queryParts, "&")
}

// 清理测试文件
func CleanupTestFile(t *testing.T, filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			t.Logf("清理测试文件失败: %v", err)
		}
	}
}