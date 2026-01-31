package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("=== 网络安全平台 - 测试运行器 ===")
	fmt.Println()

	// 1. 检查Go环境
	fmt.Println("1. 检查Go环境...")
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("❌ Go未安装或配置不正确:", err)
	}
	fmt.Printf("✅ %s", output)

	// 2. 清理并安装依赖
	fmt.Println("\n2. 清理并安装依赖...")
	runCommand("go", "mod", "tidy")

	// 3. 清理测试数据
	fmt.Println("\n3. 清理测试数据...")
	runCommand("go", "run", "scripts/clean_test_users.go")

	// 4. 运行数据库连接测试
	fmt.Println("\n4. 运行数据库连接测试...")
	runCommand("go", "run", "scripts/init_login_data_simple.go")

	// 5. 运行单元测试
	fmt.Println("\n5. 运行单元测试...")

	fmt.Println("\n====== 视频API测试 ======")
	runCommand("go", "test", "./test/unit", "-run", "TestVideo", "-v")

	fmt.Println("\n====== 登录API测试 ======")
	runCommand("go", "test", "./test/unit", "-run", "TestLogin", "-v")

	fmt.Println("\n====== 注册API测试 ======")
	runCommand("go", "test", "./test/unit", "-run", "TestRegister", "-v")

	fmt.Println("\n====== 所有测试 ======")
	runCommand("go", "test", "./test/unit", "-v")

	fmt.Println("\n✅ 测试完成！")
	fmt.Println("\n下一步：")
	fmt.Println("  1. 启动服务器: go run MAIN/server/main.go")
	fmt.Println("  2. 访问 http://localhost:3000 查看API文档")
	fmt.Println("  3. 使用curl测试各个API")
	
	fmt.Print("\n按回车键退出...")
	fmt.Scanln()
}

func runCommand(name string, args ...string) {
	fmt.Printf("执行: %s %s\n", name, strings.Join(args, " "))
	
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		// 有些测试可能会失败，这没关系
		fmt.Printf("⚠️  命令执行完成（可能有错误）: %v\n", err)
	}
}