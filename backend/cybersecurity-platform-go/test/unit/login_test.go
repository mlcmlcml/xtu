package unit

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cybersecurity-platform-go/internal/database"
	"cybersecurity-platform-go/internal/handlers"
	"golang.org/x/crypto/bcrypt"
)

// TestLoginHandler 测试登录功能
func TestLoginHandler(t *testing.T) {
	fmt.Println("开始测试登录功能...")

	// 1. 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		t.Fatalf("获取数据库连接失败: %v", err)
	}
	fmt.Println("✓ 获取数据库连接成功")

	// 2. 准备测试数据
	fmt.Println("准备测试用户数据...")
	prepareTestUserData(t, db)

	// 3. 测试用例
	testCases := []struct {
		name           string
		stuID          string
		password       string
		expectedCode   int
		expectedStatus int
		shouldSuccess  bool
	}{
		{
			name:           "正确的学号和密码",
			stuID:          "20230001",
			password:       "123456",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
			shouldSuccess:  true,
		},
		{
			name:           "正确的学号，错误的密码",
			stuID:          "20230001",
			password:       "wrongpassword",
			expectedCode:   40002,
			expectedStatus: http.StatusUnauthorized,
			shouldSuccess:  false,
		},
		{
			name:           "不存在的学号",
			stuID:          "99999999",
			password:       "123456",
			expectedCode:   40001,
			expectedStatus: http.StatusUnauthorized,
			shouldSuccess:  false,
		},
		{
			name:           "学号为空",
			stuID:          "",
			password:       "123456",
			expectedCode:   40000,
			expectedStatus: http.StatusBadRequest,
			shouldSuccess:  false,
		},
		{
			name:           "密码为空",
			stuID:          "20230001",
			password:       "",
			expectedCode:   40000,
			expectedStatus: http.StatusBadRequest,
			shouldSuccess:  false,
		},
	}

	// 4. 运行每个测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s\n", tc.name)
			fmt.Printf("学号: %s, 密码: %s\n", tc.stuID, "***")

			// 创建登录请求体
			loginReq := map[string]string{
				"stuId":    tc.stuID,
				"password": tc.password,
			}

			// 转换为JSON
			jsonData, err := json.Marshal(loginReq)
			if err != nil {
				t.Fatalf("创建JSON失败: %v", err)
			}

			// 创建HTTP请求
			req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// 创建HTTP响应记录器
			rr := httptest.NewRecorder()

			// 调用登录处理器
			handlers.LoginHandler(rr, req)

			// 检查HTTP状态码
			if rr.Code != tc.expectedStatus {
				t.Errorf("❌ HTTP状态码不正确")
				t.Errorf("  期望: %d (%s)", tc.expectedStatus, http.StatusText(tc.expectedStatus))
				t.Errorf("  实际: %d (%s)", rr.Code, http.StatusText(rr.Code))
			} else {
				fmt.Printf("✓ HTTP状态码正确: %d\n", rr.Code)
			}

			// 检查响应内容
			if tc.shouldSuccess {
				// 登录成功的情况
				var response handlers.LoginResponse
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

				// 检查消息
				if response.Message != "登录成功" {
					t.Errorf("❌ 消息不正确")
					t.Errorf("  期望: '登录成功'")
					t.Errorf("  实际: '%s'", response.Message)
				} else {
					fmt.Printf("✓ 消息正确: %s\n", response.Message)
				}

				// 检查用户信息
				if response.User == nil {
					t.Error("❌ 用户信息不应为空")
				} else {
					if response.User.StuID != tc.stuID {
						t.Errorf("❌ 学号不正确")
						t.Errorf("  期望: %s", tc.stuID)
						t.Errorf("  实际: %s", response.User.StuID)
					} else {
						fmt.Printf("✓ 学号正确: %s\n", response.User.StuID)
					}

					if response.User.ID <= 0 {
						t.Error("❌ 用户ID应为正数")
					} else {
						fmt.Printf("✓ 用户ID正确: %d\n", response.User.ID)
					}
				}

			} else {
				// 登录失败的情况
				var response handlers.LoginResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Errorf("❌ 解析JSON响应失败: %v", err)
					return
				}

				// 检查响应码
				if response.Code != tc.expectedCode {
					t.Errorf("❌ 错误码不正确")
					t.Errorf("  期望: %d", tc.expectedCode)
					t.Errorf("  实际: %d", response.Code)
				} else {
					fmt.Printf("✓ 错误码正确: %d\n", response.Code)
				}

				// 检查消息
				if response.Message == "" {
					t.Error("❌ 错误消息不应为空")
				} else {
					fmt.Printf("✓ 错误消息: %s\n", response.Message)
				}

				// 用户信息应为nil
				if response.User != nil {
					t.Error("❌ 失败时用户信息应为空")
				}
			}
		})
	}

	// 5. 测试无效的HTTP方法
	t.Run("测试无效的HTTP方法", func(t *testing.T) {
		fmt.Println("\n测试: 使用GET方法调用登录API（应该失败）")

		req := httptest.NewRequest("GET", "/api/login", nil)
		rr := httptest.NewRecorder()

		handlers.LoginHandler(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("❌ GET请求应该返回405")
			t.Errorf("  期望: 405")
			t.Errorf("  实际: %d", rr.Code)
		} else {
			fmt.Println("✓ GET请求正确返回405")
		}
	})

	fmt.Println("\n✅ 登录功能测试完成！")
}

// TestLoginRoutes 测试登录路由
func TestLoginRoutes(t *testing.T) {
	fmt.Println("\n开始测试登录路由...")

	// 创建路由
	mux := handlers.RegisterLoginRoutes()

	// 测试POST请求
	loginReq := map[string]string{
		"stuId":    "20230001",
		"password": "123456",
	}

	jsonData, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusUnauthorized {
		t.Errorf("POST请求状态码异常: %d", rr.Code)
	} else {
		fmt.Printf("✓ POST请求状态码: %d\n", rr.Code)
	}

	// 测试GET请求（应该失败）
	req = httptest.NewRequest("GET", "/api/login", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET请求应该返回405: %d", rr.Code)
	} else {
		fmt.Println("✓ GET请求正确返回405")
	}

	fmt.Println("\n✅ 登录路由测试完成！")
}

// prepareTestUserData 准备测试用户数据
func prepareTestUserData(t *testing.T, db *sql.DB) {
	t.Helper()

	fmt.Println("检查并创建用户相关表...")

	// 创建students表
	createStudentsSQL := `
		CREATE TABLE IF NOT EXISTS students (
			stuId VARCHAR(50) PRIMARY KEY,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE
		)
	`

	_, err := db.Exec(createStudentsSQL)
	if err != nil {
		t.Fatalf("创建students表失败: %v", err)
	}
	fmt.Println("✓ students表已存在或创建成功")

	// 创建userdetail表
	createUserDetailSQL := `
		CREATE TABLE IF NOT EXISTS userdetail (
			id INT PRIMARY KEY AUTO_INCREMENT,
			stuId VARCHAR(50) NOT NULL UNIQUE,
			nickName VARCHAR(50),
			userHead VARCHAR(500),
			userName VARCHAR(50),
			userEmail VARCHAR(100)
		)
	`

	_, err = db.Exec(createUserDetailSQL)
	if err != nil {
		t.Fatalf("创建userdetail表失败: %v", err)
	}
	fmt.Println("✓ userdetail表已存在或创建成功")

	// 检查是否有测试数据
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM students WHERE stuId IN ('20230001', '20230002')").Scan(&count)
	if err != nil {
		t.Logf("查询用户数量失败: %v", err)
	}

	// 如果测试数据不够，插入新数据
	if count < 2 {
		fmt.Println("插入测试用户数据...")

		// 生成bcrypt哈希密码
		password := "123456"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("生成bcrypt哈希失败: %v", err)
		}

		// 插入students表数据
		insertStudentsSQL := `
			INSERT INTO students (stuId, password, email) VALUES
			(?, ?, ?),
			(?, ?, ?)
			ON DUPLICATE KEY UPDATE
			password = VALUES(password),
			email = VALUES(email)
		`

		_, err = db.Exec(insertStudentsSQL,
			"20230001", string(hashedPassword), "student1@example.com",
			"20230002", string(hashedPassword), "student2@example.com",
		)

		if err != nil {
			t.Fatalf("插入students数据失败: %v", err)
		}

		// 插入userdetail表数据
		insertUserDetailSQL := `
			INSERT INTO userdetail (stuId, nickName, userHead, userName, userEmail) VALUES
			(?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			nickName = VALUES(nickName),
			userHead = VALUES(userHead),
			userName = VALUES(userName),
			userEmail = VALUES(userEmail)
		`

		_, err = db.Exec(insertUserDetailSQL,
			"20230001", "小明", "https://example.com/avatar1.jpg", "张三", "student1@example.com",
			"20230002", "小红", "https://example.com/avatar2.jpg", "李四", "student2@example.com",
		)

		if err != nil {
			t.Fatalf("插入userdetail数据失败: %v", err)
		}

		fmt.Println("✓ 测试用户数据插入成功")
	} else {
		fmt.Println("✓ 测试用户数据已存在")
	}

	// 显示当前用户数据
	fmt.Println("\n当前测试用户：")
	rows, err := db.Query(`
		SELECT s.stuId, u.nickName, u.userName 
		FROM students s
		LEFT JOIN userdetail u ON s.stuId = u.stuId
		WHERE s.stuId IN ('20230001', '20230002')
		ORDER BY s.stuId
	`)

	if err != nil {
		t.Logf("查询用户数据失败: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var stuID, nickName, userName string
		err := rows.Scan(&stuID, &nickName, &userName)
		if err != nil {
			t.Logf("读取用户数据失败: %v", err)
			continue
		}

		fmt.Printf("  学号: %s, 昵称: %s, 姓名: %s\n", stuID, nickName, userName)
	}
}