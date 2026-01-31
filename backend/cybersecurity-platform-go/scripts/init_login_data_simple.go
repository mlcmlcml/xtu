package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=== 初始化登录测试数据（简化版） ===")
	
	// 数据库连接
	dsn := "root:219332@tcp(localhost:3306)/cybersecurity-platform?charset=utf8mb4"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()
	
	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("✓ 数据库连接成功")
	
	// 1. 创建students表
	fmt.Println("\n1. 创建students表...")
	createStudentsTableSQL := `
		CREATE TABLE IF NOT EXISTS students (
			stuId VARCHAR(50) PRIMARY KEY,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE
		)
	`
	
	_, err = db.Exec(createStudentsTableSQL)
	if err != nil {
		log.Fatal("创建students表失败:", err)
	}
	fmt.Println("✓ students表创建成功")
	
	// 2. 创建userdetail表
	fmt.Println("\n2. 创建userdetail表...")
	createUserDetailTableSQL := `
		CREATE TABLE IF NOT EXISTS userdetail (
			id INT PRIMARY KEY AUTO_INCREMENT,
			stuId VARCHAR(50) NOT NULL UNIQUE,
			nickName VARCHAR(50),
			userHead VARCHAR(500),
			userName VARCHAR(50),
			userEmail VARCHAR(100)
		)
	`
	
	_, err = db.Exec(createUserDetailTableSQL)
	if err != nil {
		log.Fatal("创建userdetail表失败:", err)
	}
	fmt.Println("✓ userdetail表创建成功")
	
	// 3. 插入测试用户数据
	fmt.Println("\n3. 插入测试用户数据...")
	
	// 先删除可能存在的旧数据
	_, _ = db.Exec("DELETE FROM userdetail WHERE stuId IN ('20230001', '20230002')")
	_, _ = db.Exec("DELETE FROM students WHERE stuId IN ('20230001', '20230002')")
	
	// 使用Node.js项目中的bcrypt哈希值
	// 这是"123456"的bcrypt加密值：$2a$10$N9qo8uLOickgx2ZMRZoMye3Z7c3K3K9Z7mZQ7JZkFvJX9vY7qXqZC
	hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMye3Z7c3K3K9Z7mZQ7JZkFvJX9vY7qXqZC"
	
	// 插入students表数据
	insertStudentsSQL := `
		INSERT INTO students (stuId, password, email) VALUES
		(?, ?, ?),
		(?, ?, ?)
	`
	
	_, err = db.Exec(insertStudentsSQL,
		"20230001", hashedPassword, "student1@example.com",
		"20230002", hashedPassword, "student2@example.com",
	)
	
	if err != nil {
		log.Fatal("插入students数据失败:", err)
	}
	fmt.Println("✓ students表数据插入成功")
	
	// 插入userdetail表数据
	insertUserDetailSQL := `
		INSERT INTO userdetail (stuId, nickName, userHead, userName, userEmail) VALUES
		(?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?)
	`
	
	_, err = db.Exec(insertUserDetailSQL,
		"20230001", "小明", "https://example.com/avatar1.jpg", "张三", "student1@example.com",
		"20230002", "小红", "https://example.com/avatar2.jpg", "李四", "student2@example.com",
	)
	
	if err != nil {
		log.Fatal("插入userdetail数据失败:", err)
	}
	fmt.Println("✓ userdetail表数据插入成功")
	
	// 4. 验证数据
	fmt.Println("\n4. 验证数据...")
	
	query := `
		SELECT 
			s.stuId,
			s.email,
			u.nickName,
			u.userName,
			u.userHead
		FROM students s
		JOIN userdetail u ON s.stuId = u.stuId
		ORDER BY s.stuId
	`
	
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("查询数据失败:", err)
	}
	defer rows.Close()
	
	fmt.Println("┌──────────┬───────────────────────┬──────────┬──────────┬─────────────────────────────┐")
	fmt.Println("│ 学号     │ 邮箱                  │ 昵称     │ 姓名     │ 头像URL                     │")
	fmt.Println("├──────────┼───────────────────────┼──────────┼──────────┼─────────────────────────────┤")
	
	count := 0
	for rows.Next() {
		var stuID, email, nickName, userName, userHead string
		
		err := rows.Scan(&stuID, &email, &nickName, &userName, &userHead)
		if err != nil {
			log.Fatal("读取数据失败:", err)
		}
		
		// 缩短URL显示
		shortURL := userHead
		if len(shortURL) > 25 {
			shortURL = shortURL[:22] + "..."
		}
		
		fmt.Printf("│ %-8s │ %-21s │ %-8s │ %-8s │ %-27s │\n", 
			stuID, email, nickName, userName, shortURL)
		count++
	}
	fmt.Println("└──────────┴───────────────────────┴──────────┴──────────┴─────────────────────────────┘")
	
	fmt.Printf("\n✅ 登录测试数据初始化完成！共 %d 个测试用户。\n", count)
	fmt.Println("\n测试用户信息：")
	fmt.Println("   学号: 20230001, 密码: 123456")
	fmt.Println("   学号: 20230002, 密码: 123456")
	fmt.Println("\n现在可以测试登录API了！")
}