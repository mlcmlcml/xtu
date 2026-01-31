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

// TestGetVideoByID 测试获取视频详情的功能
func TestGetVideoByID(t *testing.T) {
	fmt.Println("开始测试 GetVideoByID...")
	
	// 1. 首先获取数据库连接
	// GetDB() 会返回数据库连接，如果还没有连接，它会创建连接
	db, err := database.GetDB()
	if err != nil {
		// 如果获取失败，标记测试失败
		t.Fatalf("获取数据库连接失败: %v", err)
	}
	fmt.Println("✓ 获取数据库连接成功")
	
	// 2. 准备测试数据
	// 我们先确保数据库中有测试数据
	fmt.Println("准备测试数据...")
	prepareTestVideoData(t, db)
	
	// 3. 测试各种情况
	// 我们创建三个测试用例：
	testCases := []struct {
		name           string    // 测试用例名称
		videoID        string    // 要测试的视频ID
		expectedCode   int       // 期望的响应码（20000表示成功，404表示不存在等）
		expectedStatus int       // HTTP状态码（200, 404等）
		shouldExist    bool      // 视频是否应该存在
	}{
		{
			name:           "获取存在的视频（ID=1）",
			videoID:        "1",
			expectedCode:   20000,      // 对应Node.js的成功码
			expectedStatus: http.StatusOK, // HTTP 200
			shouldExist:    true,
		},
		{
			name:           "获取不存在的视频（ID=999）",
			videoID:        "999",
			expectedCode:   404,           // 对应Node.js的404
			expectedStatus: http.StatusNotFound, // HTTP 404
			shouldExist:    false,
		},
		{
			name:           "无效的视频ID（abc不是数字）",
			videoID:        "abc",
			expectedCode:   400,                // 对应Node.js的400
			expectedStatus: http.StatusBadRequest, // HTTP 400
			shouldExist:    false,
		},
	}
	
	// 4. 运行每个测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n测试: %s\n", tc.name)
			fmt.Printf("视频ID: %s\n", tc.videoID)
			
			// 创建HTTP请求
			// 模拟浏览器访问 /api/videos/1
			req := httptest.NewRequest("GET", "/api/videos/"+tc.videoID, nil)
			
			// 创建HTTP响应记录器
			// 它会记录API返回的所有内容
			rr := httptest.NewRecorder()
			
			// 设置请求的URL路径
			// 这样GetVideoByID函数才能获取到视频ID
			req.URL.Path = "/api/videos/" + tc.videoID
			
			// 5. 调用我们写的处理函数
			handlers.GetVideoByID(rr, req)
			
			// 6. 检查HTTP状态码
			if rr.Code != tc.expectedStatus {
				t.Errorf("❌ HTTP状态码不正确")
				t.Errorf("  期望: %d (%s)", tc.expectedStatus, http.StatusText(tc.expectedStatus))
				t.Errorf("  实际: %d (%s)", rr.Code, http.StatusText(rr.Code))
			} else {
				fmt.Printf("✓ HTTP状态码正确: %d\n", rr.Code)
			}
			
			// 7. 检查响应内容
			if tc.shouldExist {
				// 视频应该存在的情况
				var response handlers.VideoResponse
				
				// 把JSON响应解析到结构体
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
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
				
				// 检查视频数据
				if response.Data.Video.ID <= 0 {
					t.Error("❌ 视频ID应该是正数")
				} else {
					fmt.Printf("✓ 视频ID正确: %d\n", response.Data.Video.ID)
				}
				
				if response.Data.Video.URL == "" {
					t.Error("❌ 视频URL不应为空")
				} else {
					fmt.Printf("✓ 视频URL存在\n")
				}
				
				if response.Data.Video.Description == "" {
					t.Error("❌ 视频描述不应为空")
				} else {
					fmt.Printf("✓ 视频描述存在\n")
				}
				
			} else {
				// 视频不存在或ID无效的情况
				var errorResp handlers.ErrorResponse
				
				err := json.NewDecoder(rr.Body).Decode(&errorResp)
				if err != nil {
					t.Errorf("❌ 解析错误响应失败: %v", err)
					return
				}
				
				if errorResp.Code != tc.expectedCode {
					t.Errorf("❌ 错误码不正确")
					t.Errorf("  期望: %d", tc.expectedCode)
					t.Errorf("  实际: %d", errorResp.Code)
				} else {
					fmt.Printf("✓ 错误码正确: %d\n", errorResp.Code)
				}
				
				if errorResp.Message == "" {
					t.Error("❌ 错误消息不应为空")
				} else {
					fmt.Printf("✓ 错误消息: %s\n", errorResp.Message)
				}
			}
		})
	}
	
	fmt.Println("\n✅ 所有测试用例完成！")
}

// TestVideoRoutes 测试视频路由
func TestVideoRoutes(t *testing.T) {
	fmt.Println("\n开始测试视频路由...")
	
	// 1. 创建路由
	// RegisterVideoRoutes() 会返回配置好的路由
	mux := handlers.RegisterVideoRoutes()
	
	// 2. 测试GET请求（应该成功）
	fmt.Println("测试有效的GET请求...")
	req := httptest.NewRequest("GET", "/api/videos/1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	
	// GET请求应该返回200（成功）或404（不存在），但不能是其他错误
	if rr.Code != http.StatusOK && rr.Code != http.StatusNotFound {
		t.Errorf("❌ GET请求状态码异常")
		t.Errorf("  期望: 200 (成功) 或 404 (不存在)")
		t.Errorf("  实际: %d", rr.Code)
	} else {
		fmt.Printf("✓ GET请求状态码: %d\n", rr.Code)
	}
	
	// 3. 测试POST请求（应该失败，因为只允许GET）
	fmt.Println("\n测试无效的POST请求...")
	req = httptest.NewRequest("POST", "/api/videos/1", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("❌ POST请求应该返回405（方法不允许）")
		t.Errorf("  期望: 405")
		t.Errorf("  实际: %d", rr.Code)
	} else {
		fmt.Printf("✓ POST请求正确返回405\n")
	}
	
	// 4. 检查错误响应
	var errorResp handlers.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errorResp)
	if err != nil {
		t.Errorf("❌ 解析错误响应失败: %v", err)
		return
	}
	
	if errorResp.Code != 405 {
		t.Errorf("❌ 错误码应该是405")
		t.Errorf("  期望: 405")
		t.Errorf("  实际: %d", errorResp.Code)
	} else {
		fmt.Printf("✓ 错误码正确: 405\n")
	}
	
	if errorResp.Message != "方法不允许" {
		t.Errorf("❌ 错误消息不正确")
		t.Errorf("  期望: '方法不允许'")
		t.Errorf("  实际: '%s'", errorResp.Message)
	} else {
		fmt.Printf("✓ 错误消息正确: %s\n", errorResp.Message)
	}
	
	fmt.Println("\n✅ 路由测试完成！")
}

// prepareTestVideoData 准备测试视频数据
func prepareTestVideoData(t *testing.T, db *sql.DB) {
	t.Helper() // 标记这是辅助函数
	
	fmt.Println("检查并创建videos表...")
	
	// 创建表的SQL语句
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS videos (
			id INT PRIMARY KEY AUTO_INCREMENT,
			url VARCHAR(500) NOT NULL,
			description TEXT,
			duration INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	
	// 执行创建表的SQL
	_, err := db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("❌ 创建videos表失败: %v", err)
	}
	fmt.Println("✓ videos表已存在或创建成功")
	
	// 检查是否有测试数据
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM videos WHERE id IN (1, 2)").Scan(&count)
	if err != nil {
		t.Logf("查询数据数量失败: %v", err)
	}
	
	// 如果测试数据不够，插入新数据
	if count < 2 {
		fmt.Println("插入测试数据...")
		
		// 插入测试数据的SQL
		insertSQL := `
			INSERT INTO videos (id, url, description, duration) VALUES
			(1, 'http://localhost:3000/api/videoing/test1.mp4', '测试视频1 - 网络安全基础', 3600),
			(2, 'http://localhost:3000/api/videoing/test2.mp4', '测试视频2 - 密码学入门', 4200)
			ON DUPLICATE KEY UPDATE
			url = VALUES(url),
			description = VALUES(description),
			duration = VALUES(duration)
		`
		
		_, err = db.Exec(insertSQL)
		if err != nil {
			t.Fatalf("❌ 插入测试数据失败: %v", err)
		}
		fmt.Println("✓ 测试数据插入成功")
	} else {
		fmt.Println("✓ 测试数据已存在")
	}
	
	// 显示当前数据
	fmt.Println("\n当前视频数据:")
	rows, err := db.Query("SELECT id, url, description, duration FROM videos ORDER BY id")
	if err != nil {
		t.Logf("查询数据失败: %v", err)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var id, duration int
		var url, description string
		
		err := rows.Scan(&id, &url, &description, &duration)
		if err != nil {
			t.Logf("读取数据失败: %v", err)
			continue
		}
		
		// 缩短显示
		shortURL := url
		if len(shortURL) > 40 {
			shortURL = shortURL[:37] + "..."
		}
		
		shortDesc := description
		if len(shortDesc) > 30 {
			shortDesc = shortDesc[:27] + "..."
		}
		
		fmt.Printf("  ID:%d 时长:%ds %s\n", id, duration, shortDesc)
	}
}