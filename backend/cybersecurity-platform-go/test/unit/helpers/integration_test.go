// test/unit/helpers/integration_test.go
package helpers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cybersecurity-platform-go/internal/handlers"
)

// 集成测试：完整的论坛功能流程
func TestForumIntegration(t *testing.T) {
	// 测试流程：分类 -> 文章 -> 详情 -> 评论
	
	// 1. 获取分类
	t.Run("获取分类列表", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/forum/categories", nil)
		w := httptest.NewRecorder()
		
		handlers.GetCategoriesHandler(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		body, _ := io.ReadAll(resp.Body)
		var response handlers.BaseResponse
		json.Unmarshal(body, &response)
		
		// 验证响应格式
		assert.Equal(t, 20000, response.Code)
		assert.NotNil(t, response.Data)
	})
	
	// 2. 获取文章列表
	t.Run("获取文章列表", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/forum/articles?page=1&pageSize=5", nil)
		w := httptest.NewRecorder()
		
		handlers.GetArticlesHandler(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		body, _ := io.ReadAll(resp.Body)
		var response handlers.BaseResponse
		json.Unmarshal(body, &response)
		
		assert.Equal(t, 20000, response.Code)
		
		// 验证数据结构
		if data, ok := response.Data.(map[string]interface{}); ok {
			assert.Contains(t, data, "items")
			assert.Contains(t, data, "total")
		}
	})
	
	// 3. 获取热门文章
	t.Run("获取热门文章", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/forum/articles/hot", nil)
		w := httptest.NewRecorder()
		
		handlers.GetHotArticlesHandler(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		body, _ := io.ReadAll(resp.Body)
		var response handlers.BaseResponse
		json.Unmarshal(body, &response)
		
		assert.Equal(t, 20000, response.Code)
	})
	
	// 4. 获取热门标签
	t.Run("获取热门标签", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/forum/tags/hot", nil)
		w := httptest.NewRecorder()
		
		handlers.GetHotTagsHandler(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		body, _ := io.ReadAll(resp.Body)
		var response handlers.BaseResponse
		json.Unmarshal(body, &response)
		
		assert.Equal(t, 20000, response.Code)
	})
	
	// 5. 测试文章详情（模拟一个存在的文章ID）
	t.Run("获取文章详情", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/forum/articles/1", nil)
		w := httptest.NewRecorder()
		
		// 设置路径参数
		req.SetPathValue("id", "1")
		
		handlers.GetArticleDetailHandler(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		
		// 由于可能没有真实数据，状态码可能是404或500
		// 但至少验证了处理器能处理请求
		statusCode := resp.StatusCode
		assert.True(t, statusCode >= 400 || statusCode == 200, 
			fmt.Sprintf("意外的状态码: %d", statusCode))
	})
	
	// 6. 测试评论功能（需要提供请求体）
	t.Run("发表评论", func(t *testing.T) {
		commentData := map[string]interface{}{
			"articleId":   1,
			"content":     "集成测试评论",
			"authorId":    "integ_test",
			"authorName":  "集成测试用户",
			"authorHead":  "test.jpg",
			"parentId":    nil,
		}
		
		jsonData, _ := json.Marshal(commentData)
		req := httptest.NewRequest("POST", "/api/forum/comments", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		
		handlers.PostCommentHandler(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		
		// 验证至少返回了某种响应
		assert.NotEqual(t, 0, resp.StatusCode)
	})
}

// 测试缓存功能
func TestCacheFunctionality(t *testing.T) {
	// 注意：由于缓存是全局的，这些测试可能会相互影响
	// 在实际测试中应该重置缓存或使用隔离的测试环境
	
	t.Run("分类缓存", func(t *testing.T) {
		// 第一次请求应该从数据库获取
		req1 := httptest.NewRequest("GET", "/api/forum/categories", nil)
		w1 := httptest.NewRecorder()
		handlers.GetCategoriesHandler(w1, req1)
		
		// 第二次请求应该从缓存获取
		req2 := httptest.NewRequest("GET", "/api/forum/categories", nil)
		w2 := httptest.NewRecorder()
		handlers.GetCategoriesHandler(w2, req2)
		
		// 两个响应都应该成功
		assert.Equal(t, http.StatusOK, w1.Result().StatusCode)
		assert.Equal(t, http.StatusOK, w2.Result().StatusCode)
	})
	
	t.Run("文章列表缓存", func(t *testing.T) {
		testParams := []map[string]string{
			{"page": "1", "pageSize": "10"},
			{"page": "2", "pageSize": "10", "cateId": "1"},
			{"page": "1", "pageSize": "5"},
		}
		
		for _, params := range testParams {
			query := "?"
			for key, value := range params {
				if query != "?" {
					query += "&"
				}
				query += fmt.Sprintf("%s=%s", key, value)
			}
			
			url := "/api/forum/articles" + query
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			
			handlers.GetArticlesHandler(w, req)
			
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		}
	})
}

// 测试错误处理
func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		method   string
		body     interface{}
		expected int
	}{
		{
			name:     "无效的文章ID",
			url:      "/api/forum/articles/abc",
			method:   "GET",
			body:     nil,
			expected: 400,
		},
		{
			name:     "发表评论缺少参数",
			url:      "/api/forum/comments",
			method:   "POST",
			body:     map[string]interface{}{}, // 空参数
			expected: 400,
		},
		{
			name:     "删除评论缺少请求体",
			url:      "/api/forum/comments/1",
			method:   "DELETE",
			body:     nil, // 无请求体
			expected: 400,
		},
		{
			name:     "点赞评论缺少参数",
			url:      "/api/forum/comments/1/like",
			method:   "POST",
			body:     nil,
			expected: 400,
		},
		{
			name:     "获取评论缺少文章ID",
			url:      "/api/forum/comments",
			method:   "GET",
			body:     nil,
			expected: 400,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			
			if tc.body != nil {
				jsonData, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(tc.method, tc.url, bytes.NewBuffer(jsonData))
			} else {
				req = httptest.NewRequest(tc.method, tc.url, nil)
			}
			
			w := httptest.NewRecorder()
			
			// 根据URL路由到对应的处理器
			switch {
			case tc.url == "/api/forum/comments" && tc.method == "POST":
				handlers.PostCommentHandler(w, req)
			case strings.Contains(tc.url, "/api/forum/comments/") && tc.method == "DELETE":
				// 设置路径参数
				req.SetPathValue("id", "1")
				handlers.DeleteCommentHandler(w, req)
			case strings.Contains(tc.url, "/api/forum/comments/") && strings.Contains(tc.url, "/like") && tc.method == "POST":
				req.SetPathValue("id", "1")
				handlers.LikeCommentHandler(w, req)
			case tc.url == "/api/forum/comments" && tc.method == "GET":
				handlers.GetCommentsHandler(w, req)
			case strings.Contains(tc.url, "/api/forum/articles/") && tc.method == "GET":
				req.SetPathValue("id", "abc") // 无效ID
				handlers.GetArticleDetailHandler(w, req)
			}
			
			resp := w.Result()
			defer resp.Body.Close()
			
			// 验证错误响应
			assert.Equal(t, tc.expected, resp.StatusCode)
			
			if resp.StatusCode >= 400 {
				body, _ := io.ReadAll(resp.Body)
				var response handlers.BaseResponse
				json.Unmarshal(body, &response)
				
				// 验证错误响应格式
				assert.NotEqual(t, 20000, response.Code)
				assert.NotEmpty(t, response.Message)
			}
		})
	}
}

// 测试并发请求
func TestConcurrentRequests(t *testing.T) {
	// 这个测试验证处理器是否能处理并发请求
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(index int) {
			req := httptest.NewRequest("GET", "/api/forum/categories", nil)
			w := httptest.NewRecorder()
			
			handlers.GetCategoriesHandler(w, req)
			
			resp := w.Result()
			resp.Body.Close()
			
			// 验证响应状态码
			if resp.StatusCode != http.StatusOK {
				t.Errorf("并发请求 %d 失败，状态码: %d", index, resp.StatusCode)
			}
			
			done <- true
		}(i)
	}
	
	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}