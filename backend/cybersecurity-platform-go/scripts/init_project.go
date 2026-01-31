// scripts/init_project.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("=== 初始化网络安全平台项目目录 ===")
	
	cwd, _ := os.Getwd()
	
	// 定义所有需要的目录
	dirs := []string{
		// assets目录结构
		filepath.Join(cwd, "assets", "forum", "articles"),
		filepath.Join(cwd, "assets", "image", "course"),
		filepath.Join(cwd, "assets", "image", "user"),
		filepath.Join(cwd, "assets", "pdf"),
		filepath.Join(cwd, "assets", "video"),
		
		// MAIN/server/static目录结构
		filepath.Join(cwd, "MAIN", "server", "static", "forum", "articles"),
		filepath.Join(cwd, "MAIN", "server", "static", "images", "user"),
		filepath.Join(cwd, "MAIN", "server", "static", "pdfs"),
		filepath.Join(cwd, "MAIN", "server", "static", "videos"),
		
		// static目录结构
		filepath.Join(cwd, "static", "forum", "articles"),
		filepath.Join(cwd, "static", "forum", "uploads"),
		filepath.Join(cwd, "static", "images", "user"),
		filepath.Join(cwd, "static", "pdfs"),
		filepath.Join(cwd, "static", "videos"),
		
		// 内部目录
		filepath.Join(cwd, "internal", "config"),
		filepath.Join(cwd, "internal", "database"),
		filepath.Join(cwd, "internal", "handlers"),
		filepath.Join(cwd, "internal", "middleware"),
		filepath.Join(cwd, "internal", "utils"),
		
		// 测试目录
		filepath.Join(cwd, "test", "unit", "helpers"),
	}
	
	successCount := 0
	failCount := 0
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("❌ 创建失败: %s (%v)\n", dir, err)
			failCount++
		} else {
			// 检查是否是新创建的
			if isNewDir(dir) {
				fmt.Printf("✅ 创建: %s\n", dir)
			} else {
				fmt.Printf("📁 已存在: %s\n", dir)
			}
			successCount++
		}
	}
	
	fmt.Printf("\n=== 完成 ===\n")
	fmt.Printf("成功: %d, 失败: %d\n", successCount, failCount)
	fmt.Printf("项目结构已准备就绪！\n")
	
	// 创建示例文件
	createSampleFiles(cwd)
}

func isNewDir(dir string) bool {
	// 简单的检查：如果目录为空，认为是新创建的
	files, _ := os.ReadDir(dir)
	return len(files) == 0
}

func createSampleFiles(cwd string) {
	fmt.Println("\n=== 创建示例文件 ===")
	
	// 创建 .env 文件
	envContent := `NODE_ENV=development
BASE_URL=http://localhost:3000
PORT=3000
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=219332
DB_NAME=cybersecurity-platform
NEO4J_URI=bolt://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=hukaile5206
`
	
	envPath := filepath.Join(cwd, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		if err := os.WriteFile(envPath, []byte(envContent), 0644); err == nil {
			fmt.Printf("✅ 创建: %s\n", envPath)
		}
	}
	
	// 创建 .env.production 文件
	envProdContent := `NODE_ENV=production
BASE_URL=http://193.112.146.64:3000
PORT=3000
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=219332
DB_NAME=cybersecurity-platform
NEO4J_URI=bolt://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=hukaile5206
`
	
	envProdPath := filepath.Join(cwd, ".env.production")
	if _, err := os.Stat(envProdPath); os.IsNotExist(err) {
		if err := os.WriteFile(envProdPath, []byte(envProdContent), 0644); err == nil {
			fmt.Printf("✅ 创建: %s\n", envProdPath)
		}
	}
	
	// 创建示例文章
	articleDir := filepath.Join(cwd, "assets", "forum", "articles")
	sampleArticle := `# 欢迎使用网络安全平台论坛

这是第一篇示例文章，用于测试论坛功能。

## 功能特点

1. 支持 Markdown 格式
2. 支持图片上传
3. 支持评论和点赞
4. 支持文章分类和标签

## 使用方法

1. 注册账号
2. 发布文章
3. 参与讨论

---

*最后更新: 2024-01-01*
`
	
	samplePath := filepath.Join(articleDir, "1.txt")
	if _, err := os.Stat(samplePath); os.IsNotExist(err) {
		if err := os.WriteFile(samplePath, []byte(sampleArticle), 0644); err == nil {
			fmt.Printf("✅ 创建示例文章: %s\n", samplePath)
		}
	}
}