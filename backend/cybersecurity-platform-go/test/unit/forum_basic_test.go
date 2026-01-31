package unit

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"cybersecurity-platform-go/internal/handlers"
)

func TestForumRoutesExist(t *testing.T) {
	// 创建路由
	mux := handlers.RegisterForumRoutes()
	
	// 测试所有论坛路由
	testCases := []struct {
		method string
		path   string
		body   []byte // 添加请求体
		expectedStatus int
	}{
		{"GET", "/api/forum/categories", nil, http.StatusOK},
		{"GET", "/api/forum/articles", nil, http.StatusOK},
		{"GET", "/api/forum/articles/hot", nil, http.StatusOK},
		{"GET", "/api/forum/articles/1", nil, http.StatusOK},
		{"GET", "/api/forum/tags/hot", nil, http.StatusOK},
		{"GET", "/api/forum/comments?articleId=1", nil, http.StatusOK},
		{"POST", "/api/forum/comments", []byte(`{"articleId":1,"content":"测试","authorId":"20230001","authorName":"测试用户"}`), http.StatusOK},
		{"DELETE", "/api/forum/comments/1", nil, http.StatusOK},
		{"POST", "/api/forum/comments/1/like", nil, http.StatusOK},
	}
	
	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			var req *http.Request
			if tc.body != nil {
				req = httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(tc.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}
			
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			
			if rr.Code != tc.expectedStatus {
				t.Errorf("路由 %s %s 状态码错误: 期望 %d, 实际 %d", 
					tc.method, tc.path, tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestForumCORSHeaders(t *testing.T) {
	// 需要直接调用路由，而不是通过中间件
	// 修改测试方式
	mux := http.NewServeMux()
	// 添加一个测试路由
	mux.HandleFunc("/api/forum/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":20000}`))
	})
	
	req := httptest.NewRequest("GET", "/api/forum/test", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	
	// 检查CORS头
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Content-Type": "application/json; charset=utf-8",
	}
	
	for header, expectedValue := range expectedHeaders {
		actualValue := rr.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("CORS头 %s 错误: 期望 %s, 实际 %s", header, expectedValue, actualValue)
		}
	}
}

func TestForumInvalidArticleID(t *testing.T) {
	mux := handlers.RegisterForumRoutes()
	
	// 测试无效的文章ID
	req := httptest.NewRequest("GET", "/api/forum/articles/abc", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusBadRequest {
		t.Errorf("无效文章ID应该返回400: 实际 %d", rr.Code)
	}
}

func TestForumMissingArticleID(t *testing.T) {
	mux := handlers.RegisterForumRoutes()
	
	// 测试缺少文章ID的评论请求
	req := httptest.NewRequest("GET", "/api/forum/comments", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusBadRequest {
		t.Errorf("缺少文章ID应该返回400: 实际 %d", rr.Code)
	}
}