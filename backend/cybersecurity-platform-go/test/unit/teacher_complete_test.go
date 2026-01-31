// test/unit/teacher_complete_test.go
package unit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cybersecurity-platform-go/internal/database"
	"cybersecurity-platform-go/internal/handlers"
)

func TestGetTeacherDetail(t *testing.T) {
	fmt.Println("开始测试获取教师详情...")

	// 确保数据库连接正常
	err := database.TestConnection()
	if err != nil {
		t.Skipf("数据库连接失败，跳过测试: %v", err)
	}

	// 测试用例
	testCases := []struct {
		name           string
		teacherID      string
		expectedCode   int
		expectedStatus int
	}{
		{
			name:           "获取存在的教师",
			teacherID:      "1",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "获取不存在的教师",
			teacherID:      "999",
			expectedCode:   404,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "无效的教师ID",
			teacherID:      "abc",
			expectedCode:   400,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s (ID=%s)\n", tc.name, tc.teacherID)

			// 创建HTTP请求
			req := httptest.NewRequest("GET", "/api/teachers/"+tc.teacherID, nil)

			// 创建HTTP响应记录器
			rr := httptest.NewRecorder()

			// 创建路由处理器
			mux := handlers.RegisterTeacherRoutes()
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
				var response handlers.TeacherDetailResponse
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

				// 检查教师信息
				teacher := response.Data.Teacher
				if teacher.ID <= 0 {
					t.Error("❌ 教师ID应为正数")
				} else {
					fmt.Printf("✓ 教师ID正确: %d\n", teacher.ID)
				}

				if teacher.TeacherName == "" {
					t.Error("❌ 教师姓名不应为空")
				} else {
					fmt.Printf("✓ 教师姓名: %s\n", teacher.TeacherName)
				}

				// 检查课程数据
				fmt.Printf("✓ 教师教授课程数: %d\n", len(response.Data.Course))
			}
		})
	}

	fmt.Println("\n✅ 教师详情功能测试完成！")
}

func TestGetTeacherList(t *testing.T) {
	fmt.Println("\n开始测试获取教师列表...")

	// 确保数据库连接正常
	err := database.TestConnection()
	if err != nil {
		t.Skipf("数据库连接失败，跳过测试: %v", err)
	}

	// 测试用例
	testCases := []struct {
		name           string
		queryParams    string
		expectedCode   int
		expectedStatus int
	}{
		{
			name:           "获取第一页教师",
			queryParams:    "?page=1&pageSize=2",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "搜索教师",
			queryParams:    "?name=教授",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "带分页和搜索",
			queryParams:    "?page=1&pageSize=5&name=老师",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s\n", tc.name)

			// 创建HTTP请求
			req := httptest.NewRequest("GET", "/api/teachers"+tc.queryParams, nil)

			// 创建HTTP响应记录器
			rr := httptest.NewRecorder()

			// 创建路由处理器
			mux := handlers.RegisterTeacherRoutes()
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
				var response handlers.TeacherListResponse
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
				fmt.Printf("✓ 获取到 %d 位教师，总数: %d\n", 
					len(response.Data.Rows), response.Data.Total)
				
				// 检查每位数师的基本信息
				for _, teacher := range response.Data.Rows {
					if teacher.ID <= 0 {
						t.Errorf("❌ 教师ID应为正数: %d", teacher.ID)
					}
					if teacher.TeacherName == "" {
						t.Error("❌ 教师姓名不应为空")
					}
				}
			}
		})
	}

	fmt.Println("\n✅ 教师列表功能测试完成！")
}