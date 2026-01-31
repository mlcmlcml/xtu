// test/unit/teacher_simple_test.go
package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"cybersecurity-platform-go/internal/handlers"
)

func TestTeacherRoutesExist(t *testing.T) {
	// 创建路由
	mux := handlers.RegisterTeacherRoutes()
	
	// 测试教师列表路由
	req := httptest.NewRequest("GET", "/api/teachers", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("教师列表路由状态码错误: 期望 %d, 实际 %d", http.StatusOK, rr.Code)
	}
	
	// 测试教师详情路由
	req = httptest.NewRequest("GET", "/api/teachers/1", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("教师详情路由状态码错误: 期望 %d, 实际 %d", http.StatusOK, rr.Code)
	}
}

func TestTeacherInvalidMethod(t *testing.T) {
	mux := handlers.RegisterTeacherRoutes()
	
	// 测试错误的方法
	req := httptest.NewRequest("POST", "/api/teachers", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("错误方法应该返回405: 期望 %d, 实际 %d", http.StatusMethodNotAllowed, rr.Code)
	}
}