package unit

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cybersecurity-platform-go/internal/database"
	"cybersecurity-platform-go/internal/handlers"
)

// TestGetCourseList 测试获取课程列表
func TestGetCourseList(t *testing.T) {
	fmt.Println("开始测试获取课程列表...")

	// 1. 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		t.Fatalf("获取数据库连接失败: %v", err)
	}
	fmt.Println("✓ 获取数据库连接成功")

	// 2. 准备测试数据
	fmt.Println("准备课程测试数据...")
	prepareCourseTestData(t, db)

	// 3. 测试用例
	testCases := []struct {
		name           string
		queryParams    string
		expectedCode   int
		expectedStatus int
	}{
		{
			name:           "获取第一页课程",
			queryParams:    "?page=1&pageSize=2",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "搜索课程",
			queryParams:    "?title=安全",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "按最新排序",
			queryParams:    "?order=1",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "按最热排序",
			queryParams:    "?order=2",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
		},
	}

	// 4. 运行每个测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s\n", tc.name)

			// 创建HTTP请求
			req := httptest.NewRequest("GET", "/api/courses"+tc.queryParams, nil)
			
			// 创建HTTP响应记录器
			rr := httptest.NewRecorder()

			// 调用课程列表处理器
			handlers.GetCourseList(rr, req)

			// 检查HTTP状态码
			if rr.Code != tc.expectedStatus {
				t.Errorf("❌ HTTP状态码不正确")
				t.Errorf("  期望: %d (%s)", tc.expectedStatus, http.StatusText(tc.expectedStatus))
				t.Errorf("  实际: %d (%s)", rr.Code, http.StatusText(rr.Code))
			} else {
				fmt.Printf("✓ HTTP状态码正确: %d\n", rr.Code)
			}

			// 检查响应内容
			var response handlers.CourseListResponse
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
			if response.Data.CourseList == nil {
				t.Error("❌ 课程列表不应为空")
			} else {
				fmt.Printf("✓ 获取到 %d 门课程\n", len(response.Data.CourseList))
				
				// 检查每门课程的基本信息
				for _, course := range response.Data.CourseList {
					if course.ID <= 0 {
						t.Errorf("❌ 课程ID应为正数: %d", course.ID)
					}
					if course.Title == "" {
						t.Error("❌ 课程标题不应为空")
					}
					if course.LessonNum < 0 {
						t.Errorf("❌ 课时数不应为负数: %d", course.LessonNum)
					}
				}
			}

			// 检查热门搜索列表
			if len(response.Data.HotList) == 0 {
				t.Error("❌ 热门搜索列表不应为空")
			} else {
				fmt.Printf("✓ 热门搜索列表: %v\n", response.Data.HotList)
			}

			// 检查科目列表（应为空数组）
			if response.Data.SubjectList == nil {
				t.Error("❌ 科目列表不应为nil")
			}
		})
	}

	// 5. 测试无效的HTTP方法
	t.Run("测试无效的HTTP方法", func(t *testing.T) {
		fmt.Println("\n测试: 使用POST方法调用课程列表API（应该失败）")

		req := httptest.NewRequest("POST", "/api/courses", nil)
		rr := httptest.NewRecorder()

		handlers.GetCourseList(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("❌ POST请求应该返回405")
			t.Errorf("  期望: 405")
			t.Errorf("  实际: %d", rr.Code)
		} else {
			fmt.Println("✓ POST请求正确返回405")
		}
	})

	fmt.Println("\n✅ 课程列表功能测试完成！")
}

// TestGetCourseDetail 测试获取课程详情
func TestGetCourseDetail(t *testing.T) {
	fmt.Println("\n开始测试获取课程详情...")

	// 1. 确保数据库连接正常（但不直接使用db）
	err := database.TestConnection()
	if err != nil {
		t.Skipf("数据库连接失败，跳过测试: %v", err)
	}

	// 2. 测试用例
	testCases := []struct {
		name           string
		courseID       string
		expectedCode   int
		expectedStatus int
		shouldExist    bool
	}{
		{
			name:           "获取存在的课程",
			courseID:       "1",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
			shouldExist:    true,
		},
		{
			name:           "获取不存在的课程",
			courseID:       "999",
			expectedCode:   404,
			expectedStatus: http.StatusNotFound,
			shouldExist:    false,
		},
		{
			name:           "无效的课程ID",
			courseID:       "abc",
			expectedCode:   400,
			expectedStatus: http.StatusBadRequest,
			shouldExist:    false,
		},
	}

	// 3. 运行每个测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s (ID=%s)\n", tc.name, tc.courseID)

			// 创建HTTP请求
			req := httptest.NewRequest("GET", "/api/courses/"+tc.courseID, nil)
			req.URL.Path = "/api/courses/" + tc.courseID

			// 创建HTTP响应记录器
			rr := httptest.NewRecorder()

			// 调用课程详情处理器
			handlers.GetCourseDetail(rr, req)

			// 检查HTTP状态码
			if rr.Code != tc.expectedStatus {
				t.Errorf("❌ HTTP状态码不正确")
				t.Errorf("  期望: %d (%s)", tc.expectedStatus, http.StatusText(tc.expectedStatus))
				t.Errorf("  实际: %d (%s)", rr.Code, http.StatusText(rr.Code))
			} else {
				fmt.Printf("✓ HTTP状态码正确: %d\n", rr.Code)
			}

			// 检查响应内容
			if tc.shouldExist {
				// 课程存在的情况
				var response handlers.CourseDetailResponse
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

				// 检查课程基本信息
				course := response.Data.Course
				if course.ID <= 0 {
					t.Error("❌ 课程ID应为正数")
				} else {
					fmt.Printf("✓ 课程ID正确: %d\n", course.ID)
				}

				if course.Title == "" {
					t.Error("❌ 课程标题不应为空")
				} else {
					fmt.Printf("✓ 课程标题: %s\n", course.Title)
				}

				// 检查章节数据
				if course.Chapter == nil {
					t.Error("❌ 章节数据不应为nil")
				} else {
					fmt.Printf("✓ 章节数量: %d\n", len(course.Chapter))
					
					// 检查每个章节
					for i, chapter := range course.Chapter {
						if chapter.Title == "" {
							t.Errorf("❌ 第 %d 章标题不应为空", i+1)
						}
						if chapter.Children == nil {
							t.Errorf("❌ 第 %d 章课时列表不应为nil", i+1)
						}
					}
				}

				// 检查教师信息
				if course.Teacher.TeacherName == "" {
					t.Error("❌ 教师姓名不应为空")
				} else {
					fmt.Printf("✓ 授课教师: %s\n", course.Teacher.TeacherName)
				}

			} else {
				// 课程不存在或ID无效的情况
				var errorResp map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&errorResp); err != nil {
					t.Errorf("❌ 解析错误响应失败: %v", err)
					return
				}

				// 检查错误码
				if code, ok := errorResp["code"].(float64); ok {
					if int(code) != tc.expectedCode {
						t.Errorf("❌ 错误码不正确")
						t.Errorf("  期望: %d", tc.expectedCode)
						t.Errorf("  实际: %d", int(code))
					} else {
						fmt.Printf("✓ 错误码正确: %d\n", int(code))
					}
				}

				// 检查错误消息
				if message, ok := errorResp["message"].(string); ok && message != "" {
					fmt.Printf("✓ 错误消息: %s\n", message)
				} else {
					t.Error("❌ 错误消息不应为空")
				}
			}
		})
	}

	fmt.Println("\n✅ 课程详情功能测试完成！")
}

// TestCourseRoutes 测试课程路由
func TestCourseRoutes(t *testing.T) {
	fmt.Println("\n开始测试课程路由...")

	// 创建路由
	mux := handlers.RegisterCourseRoutes()

	// 测试课程列表GET请求
	req := httptest.NewRequest("GET", "/api/courses?page=1&pageSize=5", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("课程列表GET请求状态码异常: %d", rr.Code)
	} else {
		fmt.Printf("✓ 课程列表GET请求状态码: %d\n", rr.Code)
	}

	// 测试课程列表POST请求（应该失败）
	req = httptest.NewRequest("POST", "/api/courses", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("课程列表POST请求应该返回405: %d", rr.Code)
	} else {
		fmt.Println("✓ 课程列表POST请求正确返回405")
	}

	// 测试课程详情GET请求
	req = httptest.NewRequest("GET", "/api/courses/1", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusNotFound {
		t.Errorf("课程详情GET请求状态码异常: %d", rr.Code)
	} else {
		fmt.Printf("✓ 课程详情GET请求状态码: %d\n", rr.Code)
	}

	fmt.Println("\n✅ 课程路由测试完成！")
}

// prepareCourseTestData 准备课程测试数据
func prepareCourseTestData(t *testing.T, db *sql.DB) {
	t.Helper()

	// 检查是否已有课程数据
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM courses").Scan(&count)
	if err != nil {
		t.Logf("查询课程数量失败: %v", err)
		return
	}

	if count == 0 {
		fmt.Println("运行课程数据初始化脚本...")
		// 这里可以调用初始化脚本，或者直接插入测试数据
		// 为了简化，我们假设数据已经通过 scripts/init_course_data.go 初始化
		t.Log("请先运行: go run scripts/init_course_data.go")
	} else {
		fmt.Printf("✓ 已有 %d 门课程数据\n", count)
	}
}