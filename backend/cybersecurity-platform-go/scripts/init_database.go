package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=== 初始化视频数据表 ===")
	
	// 数据库连接信息
	dsn := "root:219332@tcp(localhost:3306)/cybersecurity-platform?charset=utf8mb4&parseTime=True&loc=Local"
	
	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()
	
	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}
	fmt.Println("✓ 数据库连接成功")
	
	// SQL语句
	sqlStatements := []string{
		// 创建videos表
		`CREATE TABLE IF NOT EXISTS videos (
			id INT PRIMARY KEY AUTO_INCREMENT,
			url VARCHAR(500) NOT NULL,
			description TEXT,
			duration INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// 插入测试数据（如果不存在）
		`INSERT IGNORE INTO videos (id, url, description, duration) VALUES
		(1, 'http://localhost:3000/api/videoing/video1.mp4', '网络安全基础教程 - 第一部分', 3600),
		(2, 'http://localhost:3000/api/videoing/video2.mp4', '密码学原理与应用', 4200),
		(3, 'http://localhost:3000/api/videoing/video3.mp4', '网络攻防实战演示', 5400)`,
	}
	
	// 执行每个SQL语句
	for i, sql := range sqlStatements {
		fmt.Printf("\n执行SQL语句 %d/%d...\n", i+1, len(sqlStatements))
		
		_, err := db.Exec(sql)
		if err != nil {
			log.Printf("执行SQL失败: %v\n", err)
			log.Printf("SQL语句: %s\n", sql)
		} else {
			fmt.Println("✓ 执行成功")
		}
	}
	
	// 验证数据
	fmt.Println("\n=== 验证数据 ===")
	rows, err := db.Query("SELECT id, url, description, duration FROM videos ORDER BY id")
	if err != nil {
		log.Fatal("查询数据失败:", err)
	}
	defer rows.Close()
	
	var count int
	fmt.Println("ID | 时长(秒) | 描述")
	fmt.Println("---|---------|------")
	
	for rows.Next() {
		var id, duration int
		var url, description string
		
		err := rows.Scan(&id, &url, &description, &duration)
		if err != nil {
			log.Fatal("读取数据失败:", err)
		}
		
		// 缩短描述显示
		shortDesc := description
		if len(shortDesc) > 30 {
			shortDesc = shortDesc[:27] + "..."
		}
		
		fmt.Printf("%2d | %7d | %s\n", id, duration, shortDesc)
		count++
	}
	
	fmt.Printf("\n✓ 共找到 %d 条视频记录\n", count)
	
	if count == 0 {
		fmt.Println("\n⚠️  警告：没有找到视频数据")
		fmt.Println("请手动检查数据库：")
		fmt.Println("1. 运行 mysql -u root -p")
		fmt.Println("2. 输入密码：219332")
		fmt.Println("3. 执行: USE `cybersecurity-platform`;")
		fmt.Println("4. 执行: SELECT * FROM videos;")
	} else {
		fmt.Println("\n✅ 视频数据初始化完成！")
		fmt.Println("现在可以测试视频API了")
	}
}