// simple_test.go - 简单测试运行器
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "forum":
		runForumTests()
	case "all":
		runAllTests()
	case "cover":
		generateCoverage()
	case "clean":
		cleanFiles()
	case "check":
		checkEnv()
	case "help":
		printHelp()
	default:
		fmt.Println("未知命令:", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println("简单测试工具")
	fmt.Println("用法: go run simple_test.go <命令>")
	fmt.Println()
	fmt.Println("命令:")
	fmt.Println("  forum   运行论坛测试")
	fmt.Println("  all     运行所有测试")
	fmt.Println("  cover   生成覆盖率报告")
	fmt.Println("  clean   清理测试文件")
	fmt.Println("  check   检查测试环境")
	fmt.Println("  help    显示帮助")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run simple_test.go forum")
	fmt.Println("  go run simple_test.go cover")
	fmt.Println("  go run simple_test.go clean")
}

func runForumTests() {
	fmt.Println("运行论坛测试...")
	cmd := exec.Command("go", "test", "./test/unit/handlers", "-v", "-run", "Forum")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("测试失败:", err)
		os.Exit(1)
	}
}

func runAllTests() {
	fmt.Println("运行所有测试...")
	cmd := exec.Command("go", "test", "./test/unit/...", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("测试失败:", err)
		os.Exit(1)
	}
}

func generateCoverage() {
	fmt.Println("生成覆盖率报告...")
	
	// 测试
	cmd := exec.Command("go", "test", "./internal/handlers", "-coverprofile=coverage.out")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("覆盖率测试失败:", err)
		return
	}
	
	// 生成HTML
	cmd = exec.Command("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html")
	if err := cmd.Run(); err != nil {
		fmt.Println("生成HTML失败:", err)
		return
	}
	
	fmt.Println("覆盖率报告已生成:")
	fmt.Println("  coverage.out")
	fmt.Println("  coverage.html (用浏览器打开)")
}

func cleanFiles() {
	fmt.Println("清理文件...")
	files := []string{"coverage.out", "coverage.html"}
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
			fmt.Println("删除:", file)
		}
	}
	fmt.Println("完成")
}

func checkEnv() {
	fmt.Println("检查环境...")
	
	// 检查go
	cmd := exec.Command("go", "version")
	if output, err := cmd.Output(); err == nil {
		fmt.Println("Go版本:", string(output))
	} else {
		fmt.Println("Go未安装:", err)
	}
	
	// 检查目录
	dirs := []string{
		"./internal/handlers",
		"./test/unit/handlers",
		"./assets/forum/articles",
	}
	
	for _, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			fmt.Println("✓", dir, "存在")
		} else {
			fmt.Println("✗", dir, "不存在")
		}
	}
}