// scripts/integration_test.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type TestCase struct {
	Name        string
	Method      string
	URL         string
	RequestBody string
	ExpectCode  int
}

func main() {
	fmt.Println("=== 网络安全平台 - 集成测试 ===")
	fmt.Printf("测试开始时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// 读取配置文件获取服务器地址
	serverURL := "http://localhost:3000"
	if envURL := os.Getenv("TEST_SERVER_URL"); envURL != "" {
		serverURL = envURL
	}
	
	fmt.Printf("测试服务器: %s\n", serverURL)
	fmt.Println()
	
	// 定义测试用例
	testCases := []TestCase{
		// 基础API测试
		{"健康检查", "GET", "/health", "", 200},
		{"首页", "GET", "/", "", 200},
		
		// 论坛模块
		{"论坛分类", "GET", "/api/forum/categories", "", 200},
		{"热门文章", "GET", "/api/forum/articles/hot", "", 200},
		{"文章列表", "GET", "/api/forum/articles?page=1&pageSize=10", "", 200},
		{"热门标签", "GET", "/api/forum/tags/hot", "", 200},
		
		// 课程模块
		{"课程列表", "GET", "/api/courses", "", 200},
		{"课程搜索", "GET", "/api/courses?title=安全", "", 200},
		{"课程详情", "GET", "/api/courses/1", "", 200},
		
		// 教师模块
		{"教师列表", "GET", "/api/teachers", "", 200},
		{"教师详情", "GET", "/api/teachers/1", "", 200},
		
		// 图数据库
		{"初始图数据", "GET", "/api/init-graph", "", 200},
		{"扩展节点", "GET", "/api/expand-node/防火墙", "", 200},
		
		// 静态文件（HEAD请求检查）
		{"用户头像目录", "HEAD", "/img/user/", "", 200},
		{"课程图片目录", "HEAD", "/img/course/", "", 200},
		{"视频目录", "HEAD", "/api/videoing/", "", 200},
		{"PDF目录", "HEAD", "/api/pdfs/", "", 200},
	}
	
	// 运行测试
	results := runTests(serverURL, testCases)
	
	// 输出结果
	printResults(results)
	
	// 生成报告
	generateReport(results, serverURL)
}

func runTests(serverURL string, testCases []TestCase) []TestResult {
	var results []TestResult
	client := &http.Client{Timeout: 10 * time.Second}
	
	for _, tc := range testCases {
		fmt.Printf("测试: %-30s", tc.Name)
		
		result := TestResult{
			Name:   tc.Name,
			URL:    serverURL + tc.URL,
			Method: tc.Method,
		}
		
		// 创建请求
		var reqBody io.Reader
		if tc.RequestBody != "" {
			reqBody = bytes.NewBufferString(tc.RequestBody)
		}
		
		req, err := http.NewRequest(tc.Method, serverURL+tc.URL, reqBody)
		if err != nil {
			result.Status = "ERROR"
			result.Message = fmt.Sprintf("创建请求失败: %v", err)
			results = append(results, result)
			fmt.Printf(" ❌\n")
			continue
		}
		
		// 设置请求头
		if tc.Method == "POST" || tc.Method == "PUT" {
			req.Header.Set("Content-Type", "application/json")
		}
		
		// 发送请求
		start := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(start)
		
		if err != nil {
			result.Status = "ERROR"
			result.Message = fmt.Sprintf("请求失败: %v", err)
			result.Duration = elapsed
			results = append(results, result)
			fmt.Printf(" ❌\n")
			continue
		}
		
		// 读取响应体（为了计算大小）
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		result.StatusCode = resp.StatusCode
		result.Duration = elapsed
		result.ResponseSize = len(body)
		
		// 验证状态码
		if tc.ExpectCode == 0 || resp.StatusCode == tc.ExpectCode {
			result.Status = "PASS"
			fmt.Printf(" ✅ (%dms)\n", elapsed.Milliseconds())
		} else {
			result.Status = "FAIL"
			result.Message = fmt.Sprintf("期望状态码 %d，得到 %d", tc.ExpectCode, resp.StatusCode)
			fmt.Printf(" ❌ (%dms)\n", elapsed.Milliseconds())
		}
		
		results = append(results, result)
		
		// 避免请求过快
		time.Sleep(100 * time.Millisecond)
	}
	
	return results
}

type TestResult struct {
	Name         string
	URL          string
	Method       string
	Status       string // PASS, FAIL, ERROR
	StatusCode   int
	Message      string
	Duration     time.Duration
	ResponseSize int
}

func printResults(results []TestResult) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("测试结果汇总")
	fmt.Println(strings.Repeat("=", 80))
	
	passed := 0
	failed := 0
	errors := 0
	totalTime := time.Duration(0)
	
	for _, result := range results {
		totalTime += result.Duration
		
		switch result.Status {
		case "PASS":
			passed++
			fmt.Printf("✅ %-30s %-6s %-40s %6dms\n", 
				result.Name, result.Method, result.URL, result.Duration.Milliseconds())
		case "FAIL":
			failed++
			fmt.Printf("❌ %-30s %-6s %-40s %6dms - %s\n", 
				result.Name, result.Method, result.URL, result.Duration.Milliseconds(), result.Message)
		case "ERROR":
			errors++
			fmt.Printf("💥 %-30s %-6s %-40s %6dms - %s\n", 
				result.Name, result.Method, result.URL, result.Duration.Milliseconds(), result.Message)
		}
	}
	
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("总计: %d | 通过: %d | 失败: %d | 错误: %d | 总耗时: %v\n", 
		len(results), passed, failed, errors, totalTime)
	fmt.Println(strings.Repeat("=", 80))
}

func generateReport(results []TestResult, serverURL string) {
	// 创建HTML报告
	report := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>网络安全平台测试报告</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .summary { background: #f5f5f5; padding: 20px; border-radius: 5px; margin-bottom: 30px; }
        .passed { color: green; }
        .failed { color: red; }
        .error { color: orange; }
        table { width: 100%%; border-collapse: collapse; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #4CAF50; color: white; }
        tr:hover { background-color: #f5f5f5; }
    </style>
</head>
<body>
    <h1>网络安全平台测试报告</h1>
    <div class="summary">
        <h2>测试概览</h2>
        <p><strong>测试时间:</strong> %s</p>
        <p><strong>测试服务器:</strong> %s</p>
        <p><strong>测试总数:</strong> %d</p>
        <p><strong>通过:</strong> <span class="passed">%d</span> | 
           <strong>失败:</strong> <span class="failed">%d</span> | 
           <strong>错误:</strong> <span class="error">%d</span></p>
    </div>
    
    <h2>详细结果</h2>
    <table>
        <tr>
            <th>测试名称</th>
            <th>状态</th>
            <th>URL</th>
            <th>状态码</th>
            <th>响应时间</th>
            <th>响应大小</th>
            <th>消息</th>
        </tr>`, 
		time.Now().Format("2006-01-02 15:04:05"),
		serverURL,
		len(results),
		countByStatus(results, "PASS"),
		countByStatus(results, "FAIL"),
		countByStatus(results, "ERROR"))
	
	for _, result := range results {
		statusClass := strings.ToLower(result.Status)
		statusIcon := "✅"
		if result.Status == "FAIL" {
			statusIcon = "❌"
		} else if result.Status == "ERROR" {
			statusIcon = "💥"
		}
		
		report += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td class="%s">%s %s</td>
            <td><code>%s</code></td>
            <td>%d</td>
            <td>%dms</td>
            <td>%d bytes</td>
            <td>%s</td>
        </tr>`,
			result.Name,
			statusClass,
			statusIcon,
			result.Status,
			result.URL,
			result.StatusCode,
			result.Duration.Milliseconds(),
			result.ResponseSize,
			result.Message)
	}
	
	report += `
    </table>
</body>
</html>`
	
	// 保存报告
	filename := fmt.Sprintf("test_report_%s.html", time.Now().Format("20060102_150405"))
	if err := os.WriteFile(filename, []byte(report), 0644); err != nil {
		log.Printf("保存测试报告失败: %v", err)
	} else {
		fmt.Printf("\n测试报告已保存到: %s\n", filename)
	}
}

func countByStatus(results []TestResult, status string) int {
	count := 0
	for _, r := range results {
		if r.Status == status {
			count++
		}
	}
	return count
}