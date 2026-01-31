// test_forum_api.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("=== 论坛API简单测试 ===")
	fmt.Println("请确保服务器正在运行...")
	
	// 等待服务器启动
	time.Sleep(2 * time.Second)
	
	baseURL := "http://localhost:3000"
	
	// 1. 测试分类列表API
	fmt.Println("\n1. 测试分类列表API...")
	testAPI(baseURL+"/api/forum/categories", "分类列表")
	
	// 2. 测试热门标签API
	fmt.Println("\n2. 测试热门标签API...")
	testAPI(baseURL+"/api/forum/tags/hot", "热门标签")
	
	// 3. 测试其他API端点（占位功能）
	fmt.Println("\n3. 测试其他API端点...")
	testAPI(baseURL+"/api/forum/articles?page=1&pageSize=5", "文章列表")
	testAPI(baseURL+"/api/forum/articles/hot", "热门文章")
	testAPI(baseURL+"/api/forum/articles/1", "文章详情")
	testAPI(baseURL+"/api/forum/comments?articleId=1", "评论列表")
	
	fmt.Println("\n✅ 简单测试完成！")
}

func testAPI(url, name string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ %s 请求失败: %v\n", name, err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取 %s 响应失败: %v\n", name, err)
		return
	}
	
	fmt.Printf("  %s - 状态码: %d\n", name, resp.StatusCode)
	
	if resp.StatusCode == http.StatusOK {
		// 尝试解析JSON
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err == nil {
			if code, ok := data["code"].(float64); ok {
				fmt.Printf("    响应码: %.0f", code)
				if message, ok := data["message"].(string); ok && message != "" {
					fmt.Printf(", 消息: %s", message)
				}
				
				// 显示数据信息
				if dataObj, ok := data["data"].(interface{}); ok {
					switch v := dataObj.(type) {
					case []interface{}:
						fmt.Printf(", 数据条数: %d", len(v))
					case map[string]interface{}:
						if items, ok := v["items"].([]interface{}); ok {
							fmt.Printf(", 文章数: %d", len(items))
						}
						if total, ok := v["total"].(float64); ok {
							fmt.Printf(", 总数: %.0f", total)
						}
					}
				}
				fmt.Println()
			}
		}
	} else {
		fmt.Printf("    响应内容: %s\n", string(body))
	}
}