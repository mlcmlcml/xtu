package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("=== 课程API测试脚本 ===")
	fmt.Println("请确保服务器正在运行...")
	
	// 等待服务器启动
	time.Sleep(2 * time.Second)
	
	baseURL := "http://localhost:3000"
	
	// 1. 测试健康检查
	fmt.Println("\n1. 测试健康检查...")
	testHealthCheck(baseURL)
	
	// 2. 测试课程列表API
	fmt.Println("\n2. 测试课程列表API...")
	testCourseList(baseURL)
	
	// 3. 测试带参数的课程列表
	fmt.Println("\n3. 测试带参数的课程列表...")
	testCourseListWithParams(baseURL)
	
	// 4. 测试课程详情API
	fmt.Println("\n4. 测试课程详情API...")
	testCourseDetail(baseURL)
	
	fmt.Println("\n✅ 所有测试完成！")
}

func testHealthCheck(baseURL string) {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("❌ 健康检查失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		fmt.Println("✓ 健康检查通过")
	} else {
		fmt.Printf("❌ 健康检查失败，状态码: %d\n", resp.StatusCode)
	}
}

func testCourseList(baseURL string) {
	resp, err := http.Get(baseURL + "/api/courses")
	if err != nil {
		fmt.Printf("❌ 课程列表请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应长度: %d 字节\n", len(body))
	
	if resp.StatusCode == http.StatusOK {
		fmt.Println("✓ 课程列表API请求成功")
	} else {
		fmt.Printf("❌ 课程列表API请求失败\n")
		fmt.Printf("响应内容: %s\n", string(body))
	}
}

func testCourseListWithParams(baseURL string) {
	// 测试不同的查询参数
	testCases := []struct {
		name   string
		params string
	}{
		{"分页", "?page=1&pageSize=5"},
		{"搜索", "?title=安全"},
		{"热门排序", "?order=2"},
		{"综合查询", "?page=1&pageSize=3&title=网络&order=1"},
	}
	
	for _, tc := range testCases {
		url := baseURL + "/api/courses" + tc.params
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("❌ %s 请求失败: %v\n", tc.name, err)
			continue
		}
		resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("✓ %s 测试通过\n", tc.name)
		} else {
			fmt.Printf("❌ %s 测试失败，状态码: %d\n", tc.name, resp.StatusCode)
		}
	}
}

func testCourseDetail(baseURL string) {
	// 测试不同的课程ID
	courseIDs := []int{1, 2, 3, 999}
	
	for _, id := range courseIDs {
		url := fmt.Sprintf("%s/api/courses/%d", baseURL, id)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("❌ 课程%d详情请求失败: %v\n", id, err)
			continue
		}
		defer resp.Body.Close()
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("❌ 读取课程%d响应失败: %v\n", id, err)
			continue
		}
		
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("✓ 课程%d详情请求成功\n", id)
		} else if resp.StatusCode == http.StatusNotFound {
			fmt.Printf("⏭️  课程%d不存在\n", id)
		} else {
			fmt.Printf("❌ 课程%d详情请求失败，状态码: %d\n", id, resp.StatusCode)
			fmt.Printf("响应内容: %s\n", string(body))
		}
	}
}