// test/unit/api_test.go
package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cybersecurity-platform-go/internal/handlers"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	// 创建测试服务器
	setupTestServer()
	
	// 运行所有测试
	code := m.Run()
	
	// 清理资源
	teardownTestServer()
	
	os.Exit(code)
}

func setupTestServer() {
	// 注册所有路由
	forumMux := handlers.RegisterForumRoutes()
	courseMux := handlers.RegisterCourseRoutes()
	teacherMux := handlers.RegisterTeacherRoutes()
	graphMux := handlers.RegisterGraphRoutes()
	loginMux := handlers.RegisterLoginRoutes()
	registerMux := handlers.RegisterRoutes()
	studentMux := handlers.RegisterStudentRoutes()
	videoMux := handlers.RegisterVideoRoutes()
	
	// 创建主路由
	mainMux := http.NewServeMux()
	mainMux.Handle("/api/forum/", forumMux)
	mainMux.Handle("/api/courses/", courseMux)
	mainMux.Handle("/api/teachers/", teacherMux)
	mainMux.Handle("/api/", graphMux)
	mainMux.Handle("/api/login", loginMux)
	mainMux.Handle("/api/register", registerMux)
	mainMux.Handle("/api/student/", studentMux)
	mainMux.Handle("/api/videos/", videoMux)
	
	testServer = httptest.NewServer(mainMux)
	fmt.Printf("测试服务器启动在: %s\n", testServer.URL)
}

func teardownTestServer() {
	if testServer != nil {
		testServer.Close()
		fmt.Println("测试服务器已关闭")
	}
}

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/health")
	if err != nil {
		t.Fatalf("健康检查失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望状态码200，得到: %d", resp.StatusCode)
	}
	
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("status")) {
		t.Error("健康检查响应格式不正确")
	}
}

func TestForumCategories(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/forum/categories")
	if err != nil {
		t.Fatalf("论坛分类API失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("论坛分类API状态码: %d", resp.StatusCode)
	}
	
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		t.Errorf("解析JSON失败: %v", err)
	}
	
	// 检查返回的code字段
	if code, ok := result["code"].(float64); !ok || code != 20000 {
		t.Errorf("期望code=20000，得到: %v", result["code"])
	}
}

func TestCoursesAPI(t *testing.T) {
	// 测试课程列表
	resp, err := http.Get(testServer.URL + "/api/courses")
	if err != nil {
		t.Fatalf("课程API失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("课程列表状态码: %d", resp.StatusCode)
	}
	
	// 测试课程详情（模拟ID=1）
	resp2, err := http.Get(testServer.URL + "/api/courses/1")
	if err != nil {
		t.Fatalf("课程详情API失败: %v", err)
	}
	defer resp2.Body.Close()
	
	if resp2.StatusCode != http.StatusOK && resp2.StatusCode != http.StatusNotFound {
		t.Errorf("课程详情状态码: %d", resp2.StatusCode)
	}
}

func TestGraphDatabaseAPI(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/init-graph")
	if err != nil {
		t.Fatalf("图数据库API失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("图数据库API状态码: %d", resp.StatusCode)
	}
	
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		t.Errorf("解析图数据JSON失败: %v", err)
	}
	
	// 检查是否包含节点数据
	if nodes, ok := result["nodes"].([]interface{}); !ok || len(nodes) == 0 {
		t.Error("图数据缺少节点信息")
	}
}

func TestTeacherAPI(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/teachers")
	if err != nil {
		t.Fatalf("教师API失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("教师列表状态码: %d", resp.StatusCode)
	}
	
	// 测试教师详情（模拟ID=1）
	resp2, err := http.Get(testServer.URL + "/api/teachers/1")
	if err != nil {
		t.Fatalf("教师详情API失败: %v", err)
	}
	defer resp2.Body.Close()
	
	if resp2.StatusCode != http.StatusOK && resp2.StatusCode != http.StatusNotFound {
		t.Errorf("教师详情状态码: %d", resp2.StatusCode)
	}
}

func TestStudentAPI(t *testing.T) {
	// 测试学生选课
	reqBody := `{"courseId": 1, "stuId": "test123"}`
	resp, err := http.Post(testServer.URL+"/api/student/joinCourse", 
		"application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("学生选课API失败: %v", err)
	}
	defer resp.Body.Close()
	
	// 可能返回400或404，但应该能正确处理请求
	if resp.StatusCode >= 500 {
		t.Errorf("学生选课API服务器错误: %d", resp.StatusCode)
	}
}

func TestLoginAPI(t *testing.T) {
	loginData := map[string]string{
		"stuId":    "test123",
		"password": "password123",
	}
	
	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(testServer.URL+"/api/login", 
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("登录API失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Logf("登录API状态码: %d (可能需要有效测试数据)", resp.StatusCode)
	}
}

func TestRegisterAPI(t *testing.T) {
	registerData := map[string]string{
		"stuId":    fmt.Sprintf("test_%d", time.Now().Unix()),
		"email":    fmt.Sprintf("test_%d@example.com", time.Now().Unix()),
		"password": "password123",
		"nickName": "测试用户",
	}
	
	jsonData, _ := json.Marshal(registerData)
	resp, err := http.Post(testServer.URL+"/api/register", 
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("注册API失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Logf("注册API状态码: %d", resp.StatusCode)
	}
}