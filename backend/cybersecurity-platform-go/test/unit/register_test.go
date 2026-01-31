package unit

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cybersecurity-platform-go/internal/database"
	"cybersecurity-platform-go/internal/handlers"
)

// TestRegisterHandler 测试JSON注册功能
func TestRegisterHandler(t *testing.T) {
	fmt.Println("开始测试JSON注册功能...")

	// 1. 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		t.Fatalf("获取数据库连接失败: %v", err)
	}
	fmt.Println("✓ 获取数据库连接成功")

	// 2. 清理测试数据
	fmt.Println("清理旧测试数据...")
	cleanTestUsers(t, db)

	// 3. 测试用例
	testCases := []struct {
		name           string
		stuID          string
		email          string
		password       string
		nickName       string
		expectedCode   int
		expectedStatus int
		shouldSuccess  bool
		errorType      string // "学号重复", "邮箱重复", "参数错误"
	}{
		{
			name:           "正常注册",
			stuID:          "20240001",
			email:          "newuser1@example.com",
			password:       "123456",
			nickName:       "新用户1",
			expectedCode:   20000,
			expectedStatus: http.StatusOK,
			shouldSuccess:  true,
		},
		{
			name:           "学号重复",
			stuID:          "20230001", // 已存在的学号
			email:          "newuser2@example.com",
			password:       "123456",
			nickName:       "新用户2",
			expectedCode:   40001,
			expectedStatus: http.StatusBadRequest,
			shouldSuccess:  false,
			errorType:      "学号重复",
		},
		{
			name:           "邮箱重复",
			stuID:          "20240002",
			email:          "student1@example.com", // 已存在的邮箱
			password:       "123456",
			nickName:       "新用户3",
			expectedCode:   40002,
			expectedStatus: http.StatusBadRequest,
			shouldSuccess:  false,
			errorType:      "邮箱重复",
		},
		{
			name:           "学号为空",
			stuID:          "",
			email:          "newuser4@example.com",
			password:       "123456",
			nickName:       "新用户4",
			expectedCode:   40000,
			expectedStatus: http.StatusBadRequest,
			shouldSuccess:  false,
		},
		{
			name:           "邮箱为空",
			stuID:          "20240003",
			email:          "",
			password:       "123456",
			nickName:       "新用户5",
			expectedCode:   40000,
			expectedStatus: http.StatusBadRequest,
			shouldSuccess:  false,
		},
		{
			name:           "密码为空",
			stuID:          "20240004",
			email:          "newuser6@example.com",
			password:       "",
			nickName:       "新用户6",
			expectedCode:   40000,
			expectedStatus: http.StatusBadRequest,
			shouldSuccess:  false,
		},
		{
			name:           "昵称为空",
			stuID:          "20240005",
			email:          "newuser7@example.com",
			password:       "123456",
			nickName:       "",
			expectedCode:   40000,
			expectedStatus: http.StatusBadRequest,
			shouldSuccess:  false,
		},
	}

	// 4. 运行每个测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s\n", tc.name)
			fmt.Printf("学号: %s, 邮箱: %s, 昵称: %s\n", tc.stuID, tc.email, tc.nickName)

			// 创建注册请求体
			registerReq := map[string]string{
				"stuId":    tc.stuID,
				"email":    tc.email,
				"password": tc.password,
				"nickName": tc.nickName,
			}

			// 转换为JSON
			jsonData, err := json.Marshal(registerReq)
			if err != nil {
				t.Fatalf("创建JSON失败: %v", err)
			}

			// 创建HTTP请求
			req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// 创建HTTP响应记录器
			rr := httptest.NewRecorder()

			// 调用注册处理器
			handlers.RegisterHandler(rr, req)

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
				// 注册成功的情况
				var response handlers.RegisterResponse
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
				if response.Message != "注册成功" {
					t.Errorf("❌ 消息不正确")
					t.Errorf("  期望: '注册成功'")
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

					if response.User.NickName != tc.nickName {
						t.Errorf("❌ 昵称不正确")
						t.Errorf("  期望: %s", tc.nickName)
						t.Errorf("  实际: %s", response.User.NickName)
					} else {
						fmt.Printf("✓ 昵称正确: %s\n", response.User.NickName)
					}

					if response.User.UserEmail != tc.email {
						t.Errorf("❌ 邮箱不正确")
						t.Errorf("  期望: %s", tc.email)
						t.Errorf("  实际: %s", response.User.UserEmail)
					} else {
						fmt.Printf("✓ 邮箱正确: %s\n", response.User.UserEmail)
					}

					if response.User.UserHead == "" {
						t.Error("❌ 头像URL不应为空")
					} else {
						fmt.Printf("✓ 头像URL: %s\n", response.User.UserHead)
					}
				}

				// 验证数据库中的数据
				err = verifyUserInDB(t, db, tc.stuID, tc.email, tc.nickName)
				if err != nil {
					t.Errorf("❌ 数据库验证失败: %v", err)
				} else {
					fmt.Println("✓ 数据库验证成功")
				}

			} else {
				// 注册失败的情况
				var response handlers.RegisterResponse
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

				// 检查错误类型
				if tc.errorType == "学号重复" && !strings.Contains(response.Message, "学号") {
					t.Error("❌ 错误消息应包含'学号'")
				}
				if tc.errorType == "邮箱重复" && !strings.Contains(response.Message, "邮箱") {
					t.Error("❌ 错误消息应包含'邮箱'")
				}
			}
		})
	}

	// 5. 测试无效的HTTP方法
	t.Run("测试无效的HTTP方法", func(t *testing.T) {
		fmt.Println("\n测试: 使用GET方法调用注册API（应该失败）")

		req := httptest.NewRequest("GET", "/api/register", nil)
		rr := httptest.NewRecorder()

		handlers.RegisterHandler(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("❌ GET请求应该返回405")
			t.Errorf("  期望: 405")
			t.Errorf("  实际: %d", rr.Code)
		} else {
			fmt.Println("✓ GET请求正确返回405")
		}
	})

	fmt.Println("\n✅ JSON注册功能测试完成！")
}

// TestRegisterUploadHandler 测试文件上传注册功能
func TestRegisterUploadHandler(t *testing.T) {
	fmt.Println("\n开始测试文件上传注册功能...")

	// 1. 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		t.Fatalf("获取数据库连接失败: %v", err)
	}

	// 2. 清理测试数据
	cleanTestUsers(t, db)

	// 3. 创建测试图片文件
	testImagePath := createTestImage(t)
	defer os.Remove(testImagePath)

	// 4. 测试带文件上传的注册
	t.Run("带头像上传的注册", func(t *testing.T) {
		fmt.Println("\n测试: 带头像上传的注册")

		// 创建multipart表单
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 添加表单字段
		writer.WriteField("stuId", "20240010")
		writer.WriteField("email", "uploaduser@example.com")
		writer.WriteField("password", "123456")
		writer.WriteField("nickName", "上传用户")

		// 添加文件
		file, err := os.Open(testImagePath)
		if err != nil {
			t.Fatalf("打开测试图片失败: %v", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("avatar", filepath.Base(testImagePath))
		if err != nil {
			t.Fatalf("创建表单文件失败: %v", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatalf("复制文件内容失败: %v", err)
		}

		writer.Close()

		// 创建HTTP请求
		req := httptest.NewRequest("POST", "/api/register", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// 创建HTTP响应记录器
		rr := httptest.NewRecorder()

		// 调用上传注册处理器
		handlers.RegisterUploadHandler(rr, req)

		// 检查响应
		if rr.Code != http.StatusOK {
			t.Errorf("❌ 注册失败，状态码: %d", rr.Code)
			var response handlers.RegisterResponse
			if err := json.NewDecoder(rr.Body).Decode(&response); err == nil {
				t.Errorf("错误消息: %s", response.Message)
			}
		} else {
			fmt.Println("✓ 带文件上传的注册成功")
		}
	})

	// 5. 测试不带文件上传的注册（使用默认头像）
	t.Run("不带头像上传的注册", func(t *testing.T) {
		fmt.Println("\n测试: 不带头像上传的注册")

		// 创建multipart表单（不包含文件）
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 只添加表单字段
		writer.WriteField("stuId", "20240011")
		writer.WriteField("email", "noupload@example.com")
		writer.WriteField("password", "123456")
		writer.WriteField("nickName", "无上传用户")

		writer.Close()

		// 创建HTTP请求
		req := httptest.NewRequest("POST", "/api/register", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// 创建HTTP响应记录器
		rr := httptest.NewRecorder()

		// 调用上传注册处理器
		handlers.RegisterUploadHandler(rr, req)

		// 检查响应
		if rr.Code != http.StatusOK {
			t.Errorf("❌ 注册失败，状态码: %d", rr.Code)
		} else {
			fmt.Println("✓ 不带文件上传的注册成功")
			
			// 验证使用默认头像
			var response handlers.RegisterResponse
			if err := json.NewDecoder(rr.Body).Decode(&response); err == nil && response.User != nil {
				if !strings.Contains(response.User.UserHead, "bing.net") {
					fmt.Printf("✓ 使用了默认头像: %s\n", response.User.UserHead)
				}
			}
		}
	})

	// 6. 测试无效文件类型
	t.Run("测试无效文件类型", func(t *testing.T) {
		fmt.Println("\n测试: 上传非图片文件")

		// 创建txt文件
		txtPath := filepath.Join(os.TempDir(), "test.txt")
		err := os.WriteFile(txtPath, []byte("这不是图片"), 0644)
		if err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
		defer os.Remove(txtPath)

		// 创建multipart表单
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		writer.WriteField("stuId", "20240012")
		writer.WriteField("email", "invalidfile@example.com")
		writer.WriteField("password", "123456")
		writer.WriteField("nickName", "无效文件用户")

		// 添加txt文件
		file, err := os.Open(txtPath)
		if err != nil {
			t.Fatalf("打开测试文件失败: %v", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("avatar", "test.txt")
		if err != nil {
			t.Fatalf("创建表单文件失败: %v", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatalf("复制文件内容失败: %v", err)
		}

		writer.Close()

		// 创建HTTP请求
		req := httptest.NewRequest("POST", "/api/register", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// 创建HTTP响应记录器
		rr := httptest.NewRecorder()

			// 调用上传注册处理器
		handlers.RegisterUploadHandler(rr, req)

		// 应该返回错误
		if rr.Code != http.StatusBadRequest {
			t.Errorf("❌ 应拒绝非图片文件，状态码: %d", rr.Code)
		} else {
			fmt.Println("✓ 正确拒绝了非图片文件")
		}
	})

	fmt.Println("\n✅ 文件上传注册功能测试完成！")
}

// TestRegisterRoutes 测试注册路由
func TestRegisterRoutes(t *testing.T) {
	fmt.Println("\n开始测试注册路由...")

	// 创建路由
	mux := handlers.RegisterRoutes()

	// 测试JSON注册
	registerReq := map[string]string{
		"stuId":    "20240020",
		"email":    "routeuser@example.com",
		"password": "123456",
		"nickName": "路由测试用户",
	}

	jsonData, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusBadRequest {
		t.Errorf("JSON注册请求状态码异常: %d", rr.Code)
	} else {
		fmt.Printf("✓ JSON注册请求状态码: %d\n", rr.Code)
	}

	// 测试GET请求（应该失败）
	req = httptest.NewRequest("GET", "/api/register", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET请求应该返回405: %d", rr.Code)
	} else {
		fmt.Println("✓ GET请求正确返回405")
	}

	fmt.Println("\n✅ 注册路由测试完成！")
}

// cleanTestUsers 清理测试用户数据
func cleanTestUsers(t *testing.T, db *sql.DB) {
	t.Helper()

	// 删除非测试用户（保留20230001和20230002）
	_, err := db.Exec(`
		DELETE u FROM userdetail u
		JOIN students s ON u.stuId = s.stuId
		WHERE s.stuId NOT IN ('20230001', '20230002') AND s.stuId LIKE '2024%'
	`)
	if err != nil {
		t.Logf("清理用户数据失败: %v", err)
	}

	_, err = db.Exec("DELETE FROM students WHERE stuId NOT IN ('20230001', '20230002') AND stuId LIKE '2024%'")
	if err != nil {
		t.Logf("清理学生数据失败: %v", err)
	}
}

// verifyUserInDB 验证用户是否在数据库中
func verifyUserInDB(t *testing.T, db *sql.DB, stuID, email, nickName string) error {
	t.Helper()

	// 检查students表
	var dbStuID, dbEmail string
	err := db.QueryRow("SELECT stuId, email FROM students WHERE stuId = ?", stuID).Scan(&dbStuID, &dbEmail)
	if err != nil {
		return fmt.Errorf("查询students表失败: %v", err)
	}

	if dbStuID != stuID {
		return fmt.Errorf("学号不匹配: 期望 %s, 实际 %s", stuID, dbStuID)
	}

	if dbEmail != email {
		return fmt.Errorf("邮箱不匹配: 期望 %s, 实际 %s", email, dbEmail)
	}

	// 检查userdetail表
	var dbNickName string
	err = db.QueryRow("SELECT nickName FROM userdetail WHERE stuId = ?", stuID).Scan(&dbNickName)
	if err != nil {
		return fmt.Errorf("查询userdetail表失败: %v", err)
	}

	if dbNickName != nickName {
		return fmt.Errorf("昵称不匹配: 期望 %s, 实际 %s", nickName, dbNickName)
	}

	return nil
}

// createTestImage 创建测试图片文件
func createTestImage(t *testing.T) string {
	t.Helper()

	// 创建一个简单的1x1像素的PNG图片
	pngHeader := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG签名
		0x00, 0x00, 0x00, 0x0D, // IHDR块长度
		0x49, 0x48, 0x44, 0x52, // "IHDR"
		0x00, 0x00, 0x00, 0x01, // 宽度
		0x00, 0x00, 0x00, 0x01, // 高度
		0x08, 0x02, 0x00, 0x00, 0x00, // 位深、颜色类型等
		0x90, 0x77, 0x53, 0xDE, // CRC
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test*.png")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer tmpFile.Close()

	// 写入PNG数据
	_, err = tmpFile.Write(pngHeader)
	if err != nil {
		t.Fatalf("写入测试图片失败: %v", err)
	}

	return tmpFile.Name()
}