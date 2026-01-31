// test/unit/helpers/run_tests.go
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// 定义命令行参数
	testType := flag.String("type", "all", "测试类型: unit, integration, all")
	verbose := flag.Bool("v", false, "显示详细输出")
	timeout := flag.Duration("timeout", 5*time.Minute, "测试超时时间")
	flag.Parse()
	
	fmt.Println("=== 论坛功能测试开始 ===")
	fmt.Printf("测试类型: %s\n", *testType)
	fmt.Printf("超时时间: %v\n", *timeout)
	fmt.Println()
	
	// 构建测试命令
	var testArgs []string
	testArgs = append(testArgs, "test")
	
	if *verbose {
		testArgs = append(testArgs, "-v")
	}
	
	testArgs = append(testArgs, "-timeout", timeout.String())
	
	// 根据测试类型添加不同的参数
	switch strings.ToLower(*testType) {
	case "unit":
		testArgs = append(testArgs, "-run", "^Test.*Handler$")
	case "integration":
		testArgs = append(testArgs, "-run", "^TestForumIntegration$")
	case "all":
		// 运行所有测试
	default:
		fmt.Printf("未知的测试类型: %s\n", *testType)
		os.Exit(1)
	}
	
	// 添加测试目录
	testArgs = append(testArgs, "./test/unit/helpers")
	
	// 执行测试
	fmt.Printf("执行命令: go %s\n", strings.Join(testArgs, " "))
	fmt.Println()
	
	cmd := exec.Command("go", testArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)
	
	fmt.Println()
	fmt.Println("=== 测试统计 ===")
	fmt.Printf("运行时间: %v\n", duration)
	fmt.Printf("结束时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	
	if err != nil {
		fmt.Printf("\n❌ 测试失败: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("\n✅ 所有测试通过！")
	}
}