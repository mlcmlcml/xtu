package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== API测试脚本 ===")
	fmt.Println("注意：请确保服务器正在运行（go run MAIN/server/main.go）")
	fmt.Println()

	// 配置
	port := "3000"
	baseURL := "http://localhost:" + port

	// 检查服务器是否运行
	fmt.Print("检查服务器是否运行...")
	if !checkServerRunning(baseURL) {
		fmt.Println("\n❌ 服务器未运行！")
		fmt.Println("请先运行: go run MAIN/server/main.go")
		fmt.Print("\n按回车键退出...")
		fmt.Scanln()
		return
	}
	fmt.Println(" ✅")

	// 1. 测试健康检查
	fmt.Println("\n1. 测试健康检查...")
	testEndpoint(baseURL+"/health", "GET", "", "")

	// 2. 测试视频API
	fmt.Println("\n2. 测试视频API...")
	fmt.Println("获取视频ID=1:")
	testEndpoint(baseURL+"/api/videos/1", "GET", "", "")
	fmt.Println("\n获取不存在的视频:")
	testEndpoint(baseURL+"/api/videos/999", "GET", "", "")

	// 3. 测试登录API
	fmt.Println("\n3. 测试登录API...")
	fmt.Println("正确登录:")
	loginData := `{"stuId": "20230001", "password": "123456"}`
	testEndpoint(baseURL+"/api/login", "POST", "application/json", loginData)

	fmt.Println("\n密码错误:")
	wrongData := `{"stuId": "20230001", "password": "wrong"}`
	testEndpoint(baseURL+"/api/login", "POST", "application/json", wrongData)

	// 4. 测试注册API
	fmt.Println("\n4. 测试注册API...")
	fmt.Println("JSON注册:")
	registerData := `{"stuId": "20240099", "email": "test99@example.com", "password": "123456", "nickName": "测试用户"}`
	testEndpoint(baseURL+"/api/register", "POST", "application/json", registerData)

	fmt.Println("\n学号重复:")
	duplicateData := `{"stuId": "20230001", "email": "new@example.com", "password": "123456", "nickName": "重复用户"}`
	testEndpoint(baseURL+"/api/register", "POST", "application/json", duplicateData)

	fmt.Println("\n✅ API测试完成！")
	fmt.Print("\n按回车键退出...")
	fmt.Scanln()
}

// checkServerRunning 检查服务器是否在运行
func checkServerRunning(baseURL string) bool {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// testEndpoint 测试API端点
func testEndpoint(url, method, contentType, body string) {
	var req *http.Request
	var err error

	// 创建请求
	if body != "" {
		req, err = http.NewRequest(method, url, bytes.NewBufferString(body))
		if err != nil {
			fmt.Printf("   ❌ 创建请求失败: %v\n", err)
			return
		}
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			fmt.Printf("   ❌ 创建请求失败: %v\n", err)
			return
		}
	}

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("   ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("   ❌ 读取响应失败: %v\n", err)
		return
	}

	// 解析JSON为美化格式
	prettyResponse := formatJSON(string(responseBody))

	// 输出结果
	statusColor := "✅"
	if resp.StatusCode >= 400 {
		statusColor = "❌"
	}

	fmt.Printf("   %s 状态码: %d\n", statusColor, resp.StatusCode)
	if prettyResponse != "" {
		lines := strings.Split(prettyResponse, "\n")
		for _, line := range lines {
			if line != "" {
				fmt.Printf("      %s\n", line)
			}
		}
	}
}

// formatJSON 格式化JSON输出
func formatJSON(input string) string {
	// 简单美化：添加缩进
	var result strings.Builder
	depth := 0
	inString := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch ch {
		case '{', '[':
			if !inString {
				result.WriteByte(ch)
				result.WriteByte('\n')
				depth++
				result.WriteString(strings.Repeat("  ", depth))
			} else {
				result.WriteByte(ch)
			}
		case '}', ']':
			if !inString {
				result.WriteByte('\n')
				depth--
				result.WriteString(strings.Repeat("  ", depth))
				result.WriteByte(ch)
			} else {
				result.WriteByte(ch)
			}
		case ',':
			if !inString {
				result.WriteByte(ch)
				result.WriteByte('\n')
				result.WriteString(strings.Repeat("  ", depth))
			} else {
				result.WriteByte(ch)
			}
		case '"':
			if i > 0 && input[i-1] != '\\' {
				inString = !inString
			}
			result.WriteByte(ch)
		case ':':
			if !inString {
				result.WriteString(": ")
			} else {
				result.WriteByte(ch)
			}
		default:
			result.WriteByte(ch)
		}
	}

	return result.String()
}