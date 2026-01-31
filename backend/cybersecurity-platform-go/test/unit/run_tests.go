// run_tests.go - 命令行测试工具
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// 定义命令行参数
	var (
		runAll     = flag.Bool("all", false, "运行所有测试")
		runForum   = flag.Bool("forum", false, "运行论坛测试")
		runCover   = flag.Bool("cover", false, "生成覆盖率报告")
		runClean   = flag.Bool("clean", false, "清理测试文件")
		runCheck   = flag.Bool("check", false, "检查测试环境")
		verbose    = flag.Bool("v", false, "详细输出")
		showHelp   = flag.Bool("help", false, "显示帮助")
	)

	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	// 获取项目根目录
	projectRoot, err := getProjectRoot()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("项目: %s\n", filepath.Base(projectRoot))
	fmt.Printf("路径: %s\n\n", projectRoot)

	// 执行操作
	anyOperation := false

	if *runCheck {
		anyOperation = true
		checkEnvironment(projectRoot)
	}

	if *runClean {
		anyOperation = true
		cleanFiles(projectRoot)
	}

	if *runForum || (!anyOperation && !*runAll && !*runCover && !*runClean && !*runCheck) {
		anyOperation = true
		runForumTests(projectRoot, *verbose)
	}

	if *runAll {
		anyOperation = true
		runAllTests(projectRoot, *verbose)
	}

	if *runCover {
		anyOperation = true
		generateCoverageReport(projectRoot)
	}

	if !anyOperation {
		printHelp()
	}
}

func printHelp() {
	fmt.Println("论坛测试工具")
	fmt.Println("用法: go run run_tests.go [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -all    运行所有测试")
	fmt.Println("  -forum  运行论坛测试 (默认)")
	fmt.Println("  -cover  生成覆盖率报告")
	fmt.Println("  -clean  清理测试文件")
	fmt.Println("  -check  检查测试环境")
	fmt.Println("  -v      详细输出")
	fmt.Println("  -help   显示帮助")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run run_tests.go -forum -v")
	fmt.Println("  go run run_tests.go -cover")
	fmt.Println("  go run run_tests.go -clean")
}

func getProjectRoot() (string, error) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 在当前目录或父目录查找 go.mod
	current := wd
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("未找到 go.mod，请确保在项目目录中运行")
}

func runForumTests(projectRoot string, verbose bool) {
	fmt.Println("▶ 运行论坛测试...")

	args := []string{"test", "./test/unit/handlers", "-run", "Forum"}
	if verbose {
		args = append(args, "-v")
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("✗ 论坛测试失败")
		os.Exit(1)
	}
	fmt.Println("✓ 论坛测试通过")
}

func runAllTests(projectRoot string, verbose bool) {
	fmt.Println("▶ 运行所有测试...")

	args := []string{"test", "./test/unit/..."}
	if verbose {
		args = append(args, "-v")
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("✗ 测试失败")
		os.Exit(1)
	}
	fmt.Println("✓ 所有测试通过")
}

func generateCoverageReport(projectRoot string) {
	fmt.Println("▶ 生成覆盖率报告...")

	// 生成覆盖率数据
	cmd := exec.Command("go", "test", "./internal/handlers", "-coverprofile=coverage.out")
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("✗ 覆盖率测试失败")
		return
	}

	// 生成HTML报告
	cmd = exec.Command("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html")
	cmd.Dir = projectRoot
	if err := cmd.Run(); err != nil {
		fmt.Println("✗ 生成HTML报告失败")
		return
	}

	// 显示覆盖率摘要
	fmt.Println("覆盖率摘要:")
	cmd = exec.Command("go", "tool", "cover", "-func=coverage.out")
	cmd.Dir = projectRoot
	output, _ := cmd.Output()
	fmt.Println(string(output))

	fmt.Println("✓ 报告已生成:")
	fmt.Printf("  coverage.out   - 覆盖率数据\n")
	fmt.Printf("  coverage.html  - HTML报告(用浏览器打开查看)\n")
}

func cleanFiles(projectRoot string) {
	fmt.Println("▶ 清理测试文件...")

	files := []string{
		filepath.Join(projectRoot, "coverage.out"),
		filepath.Join(projectRoot, "coverage.html"),
	}

	cleaned := 0
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			if err := os.Remove(file); err == nil {
				fmt.Printf("  删除: %s\n", filepath.Base(file))
				cleaned++
			}
		}
	}

	if cleaned > 0 {
		fmt.Printf("✓ 清理了 %d 个文件\n", cleaned)
	} else {
		fmt.Println("✓ 没有需要清理的文件")
	}
}

func checkEnvironment(projectRoot string) {
	fmt.Println("▶ 检查测试环境...")

	// 检查目录
	checkDir := func(path, name string) bool {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("  ✓ %s: %s\n", name, path)
			return true
		}
		fmt.Printf("  ✗ %s: 不存在\n", name)
		return false
	}

	allOK := true

	// 检查关键目录
	allOK = checkDir(filepath.Join(projectRoot, "internal", "handlers"), "处理器目录") && allOK
	allOK = checkDir(filepath.Join(projectRoot, "test", "unit", "handlers"), "测试目录") && allOK

	// 检查文章目录
	articleDirs := []string{
		filepath.Join(projectRoot, "assets", "forum", "articles"),
		filepath.Join(projectRoot, "MAIN", "server", "static", "forum", "articles"),
		filepath.Join(projectRoot, "static", "forum", "articles"),
	}

	for _, dir := range articleDirs {
		if _, err := os.Stat(dir); err == nil {
			files, _ := os.ReadDir(dir)
			txtCount := 0
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".txt") {
					txtCount++
				}
			}
			fmt.Printf("  ✓ 文章目录: %s (%d 个txt文件)\n", filepath.Base(filepath.Dir(filepath.Dir(dir))), txtCount)
		} else {
			fmt.Printf("  ⚠ 文章目录: %s (不存在)\n", filepath.Base(filepath.Dir(filepath.Dir(dir))))
		}
	}

	// 检查Go版本
	cmd := exec.Command("go", "version")
	if output, err := cmd.Output(); err == nil {
		fmt.Printf("  ✓ %s", strings.TrimSpace(string(output)))
	} else {
		fmt.Printf("  ✗ Go未安装\n")
		allOK = false
	}

	if allOK {
		fmt.Println("\n✓ 环境检查通过")
	} else {
		fmt.Println("\n✗ 环境存在问题")
	}
}