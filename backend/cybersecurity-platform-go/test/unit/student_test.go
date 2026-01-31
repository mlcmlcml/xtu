// test/unit/student_test.go
package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cybersecurity-platform-go/internal/database"
	"cybersecurity-platform-go/internal/handlers"
)

// TestJoinCourse 测试加入课程
func TestJoinCourse(t *testing.T) {
	fmt.Println("开始测试加入课程...")

	// 确保数据库连接正常
	err := database.TestConnection()
	if err != nil {
		t.Skipf("数据库连接失败，跳过测试: %v", err)
	}

	// 测试用例
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedCode   int
		expectedStatus int
		description    string
	}{
		{
			name: "正常加入课程",
			requestBody: map[string]interface{}{
				"courseId": 1,
				"stuId":    "20230001",
			},
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
			description:    "应该成功加入课程",
		},
		{
			name: "参数不完整",
			requestBody: map[string]interface{}{
				"courseId": 1,
				// 缺少 stuId
			},
			expectedCode:   40000,
			expectedStatus: http.StatusBadRequest,
			description:    "缺少必要参数应该失败",
		},
		{
			name: "课程不存在",
			requestBody: map[string]interface{}{
				"courseId": 99999,
				"stuId":    "20230001",
			},
			expectedCode:   40400,
			expectedStatus: http.StatusNotFound,
			description:    "不存在的课程应该失败",
		},
		{
			name: "学生不存在",
			requestBody: map[string]interface{}{
				"courseId": 1,
				"stuId":    "99999999",
			},
			expectedCode:   40401,
			expectedStatus: http.StatusNotFound,
			description:    "不存在的学生应该失败",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s\n", tc.name)

			// 准备请求体
			jsonBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("创建请求体失败: %v", err)
			}

			// 创建HTTP请求
			req := httptest.NewRequest("POST", "/api/student/joinCourse", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// 创建HTTP响应记录器
			rr := httptest.NewRecorder()

			// 创建路由处理器
			mux := handlers.RegisterStudentRoutes()
			mux.ServeHTTP(rr, req)

			// 检查HTTP状态码
			if rr.Code != tc.expectedStatus {
				t.Errorf("❌ HTTP状态码不正确")
				t.Errorf("  期望: %d (%s)", tc.expectedStatus, http.StatusText(tc.expectedStatus))
				t.Errorf("  实际: %d (%s)", rr.Code, http.StatusText(rr.Code))
			} else {
				fmt.Printf("✓ HTTP状态码正确: %d\n", rr.Code)
			}

			// 检查响应内容
			var response map[string]interface{}
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Errorf("❌ 解析JSON响应失败: %v", err)
				return
			}

			// 检查响应码
			if code, ok := response["code"].(float64); ok {
				if int(code) != tc.expectedCode {
					t.Errorf("❌ 响应码不正确")
					t.Errorf("  期望: %d", tc.expectedCode)
					t.Errorf("  实际: %d", int(code))
				} else {
					fmt.Printf("✓ 响应码正确: %d\n", int(code))
				}
			}

			// 检查错误消息（如果有）
			if message, ok := response["message"].(string); ok && message != "" {
				fmt.Printf("✓ 响应消息: %s\n", message)
			}
		})
	}

	fmt.Println("\n✅ 加入课程测试完成！")
}

// TestMyCourses 测试获取我的课程
func TestMyCourses(t *testing.T) {
	fmt.Println("\n开始测试获取我的课程...")

	// 确保数据库连接正常
	err := database.TestConnection()
	if err != nil {
		t.Skipf("数据库连接失败，跳过测试: %v", err)
	}

	// 测试用例
	testCases := []struct {
		name           string
		stuID          string
		page           string
		size           string
		expectedCode   int
		expectedStatus int
	}{
		{
			name:           "正常获取",
			stuID:          "20230001",
			page:           "1",
			size:           "10",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "缺少学生ID",
			stuID:          "",
			page:           "1",
			size:           "10",
			expectedCode:   40000,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "分页参数",
			stuID:          "20230001",
			page:           "2",
			size:           "5",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s\n", tc.name)

			// 构建URL
			url := fmt.Sprintf("/api/student/myCourses?stuId=%s", tc.stuID)
			if tc.page != "" {
				url += fmt.Sprintf("&page=%s", tc.page)
			}
			if tc.size != "" {
				url += fmt.Sprintf("&size=%s", tc.size)
			}

			// 创建HTTP请求
			req := httptest.NewRequest("GET", url, nil)

			// 创建HTTP响应记录器
			rr := httptest.NewRecorder()

			// 创建路由处理器
			mux := handlers.RegisterStudentRoutes()
			mux.ServeHTTP(rr, req)

			// 检查HTTP状态码
			if rr.Code != tc.expectedStatus {
				t.Errorf("❌ HTTP状态码不正确")
				t.Errorf("  期望: %d (%s)", tc.expectedStatus, http.StatusText(tc.expectedStatus))
				t.Errorf("  实际: %d (%s)", rr.Code, http.StatusText(rr.Code))
			} else {
				fmt.Printf("✓ HTTP状态码正确: %d\n", rr.Code)
			}

			if rr.Code == http.StatusOK {
				// 检查成功响应
				var response handlers.MyCoursesResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Errorf("❌ 解析JSON响应失败: %v", err)
					return
				}

				// 检查响应码
				if response.Code != tc.expectedCode {
					t.Errorf("❌ 响应码不正确")
					t.Errorf("  期望: %d", tc.expectedCode)
					t.Errorf("  实际: %d", response.Code)
				} else {
					fmt.Printf("✓ 响应码正确: %d\n", response.Code)
				}

				// 检查数据结构
				fmt.Printf("✓ 获取到 %d 门课程，总数: %d\n", 
					len(response.Data.Courses), response.Data.Total)
			}
		})
	}

	fmt.Println("\n✅ 我的课程测试完成！")
}

// TestCheckEnrollment 测试检查选课状态
func TestCheckEnrollment(t *testing.T) {
	fmt.Println("\n开始测试检查选课状态...")

	// 确保数据库连接正常
	err := database.TestConnection()
	if err != nil {
		t.Skipf("数据库连接失败，跳过测试: %v", err)
	}

	// 测试用例
	testCases := []struct {
		name           string
		courseID       string
		stuID          string
		expectedCode   int
		expectedStatus int
		shouldEnrolled bool
	}{
		{
			name:           "已选课",
			courseID:       "1",
			stuID:          "20230001",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
			shouldEnrolled: true,
		},
		{
			name:           "未选课",
			courseID:       "3", // 假设课程3未被选
			stuID:          "20230001",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
			shouldEnrolled: false,
		},
		{
			name:           "参数不完整",
			courseID:       "",
			stuID:          "20230001",
			expectedCode:   40000,
			expectedStatus: http.StatusBadRequest,
			shouldEnrolled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s\n", tc.name)

			// 构建URL
			url := fmt.Sprintf("/api/student/checkEnrollment?courseId=%s&stuId=%s", tc.courseID, tc.stuID)

			// 创建HTTP请求
			req := httptest.NewRequest("GET", url, nil)

			// 创建HTTP响应记录器
			rr := httptest.NewRecorder()

			// 创建路由处理器
			mux := handlers.RegisterStudentRoutes()
			mux.ServeHTTP(rr, req)

			// 检查HTTP状态码
			if rr.Code != tc.expectedStatus {
				t.Errorf("❌ HTTP状态码不正确")
				t.Errorf("  期望: %d (%s)", tc.expectedStatus, http.StatusText(tc.expectedStatus))
				t.Errorf("  实际: %d (%s)", rr.Code, http.StatusText(rr.Code))
			} else {
				fmt.Printf("✓ HTTP状态码正确: %d\n", rr.Code)
			}

			if rr.Code == http.StatusOK {
				// 检查成功响应
				var response handlers.CheckEnrollmentResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Errorf("❌ 解析JSON响应失败: %v", err)
					return
				}

				// 检查响应码
				if response.Code != tc.expectedCode {
					t.Errorf("❌ 响应码不正确")
					t.Errorf("  期望: %d", tc.expectedCode)
					t.Errorf("  实际: %d", response.Code)
				} else {
					fmt.Printf("✓ 响应码正确: %d\n", response.Code)
				}

				// 检查选课状态
				fmt.Printf("✓ 选课状态: %v\n", response.Data.IsEnrolled)
			}
		})
	}

	fmt.Println("\n✅ 检查选课状态测试完成！")
}