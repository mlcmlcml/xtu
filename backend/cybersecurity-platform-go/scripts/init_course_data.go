package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=== 初始化课程测试数据 ===")
	
	dsn := "root:219332@tcp(localhost:3306)/cybersecurity-platform?charset=utf8mb4"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()
	
	err = db.Ping()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("✓ 数据库连接成功")
	
	// 创建课程相关表
	fmt.Println("\n1. 创建课程相关表...")
	createTables(db)
	
	// 插入测试数据
	fmt.Println("\n2. 插入测试数据...")
	insertTestData(db)
	
	// 验证数据
	fmt.Println("\n3. 验证数据...")
	verifyData(db)
	
	fmt.Println("\n✅ 课程测试数据初始化完成！")
}

func createTables(db *sql.DB) {
	tables := []struct {
		name string
		sql  string
	}{
		{
			name: "courses",
			sql: `
				CREATE TABLE IF NOT EXISTS courses (
					id INT PRIMARY KEY AUTO_INCREMENT,
					title VARCHAR(100) NOT NULL,
					description TEXT,
					cover VARCHAR(500),
					lesson_num INT DEFAULT 0,
					credit DECIMAL(3,1) DEFAULT 0,
					limit_count INT DEFAULT 100
				)
			`,
		},
		{
			name: "teachers",
			sql: `
				CREATE TABLE IF NOT EXISTS teachers (
					id INT PRIMARY KEY AUTO_INCREMENT,
					name VARCHAR(50) NOT NULL,
					career VARCHAR(100),
					intro TEXT
				)
			`,
		},
		{
			name: "teacher_courses",
			sql: `
				CREATE TABLE IF NOT EXISTS teacher_courses (
					id INT PRIMARY KEY AUTO_INCREMENT,
					teacher_id INT,
					course_id INT,
					FOREIGN KEY (teacher_id) REFERENCES teachers(id),
					FOREIGN KEY (course_id) REFERENCES courses(id)
				)
			`,
		},
		{
			name: "chapters",
			sql: `
				CREATE TABLE IF NOT EXISTS chapters (
					id INT PRIMARY KEY AUTO_INCREMENT,
					course_id INT,
					title VARCHAR(100) NOT NULL,
					state INT DEFAULT 0,
					FOREIGN KEY (course_id) REFERENCES courses(id)
				)
			`,
		},
		{
			name: "chapter_children",
			sql: `
				CREATE TABLE IF NOT EXISTS chapter_children (
					id INT PRIMARY KEY AUTO_INCREMENT,
					chapter_id INT,
					title VARCHAR(100) NOT NULL,
					video_id INT,
					FOREIGN KEY (chapter_id) REFERENCES chapters(id)
				)
			`,
		},
		{
			name: "videos",
			sql: `
				CREATE TABLE IF NOT EXISTS videos (
					id INT PRIMARY KEY AUTO_INCREMENT,
					url VARCHAR(500) NOT NULL,
					description TEXT,
					duration INT DEFAULT 0
				)
			`,
		},
	}
	
	for _, table := range tables {
		_, err := db.Exec(table.sql)
		if err != nil {
			log.Printf("创建 %s 表失败: %v", table.name, err)
		} else {
			fmt.Printf("✓ %s 表创建成功\n", table.name)
		}
	}
}

func insertTestData(db *sql.DB) {
	// 清空旧数据
	tables := []string{
		"teacher_courses",
		"chapter_children", 
		"chapters",
		"videos",
		"teacher_courses",
		"teachers",
		"courses",
	}
	
	for _, table := range tables {
		_, _ = db.Exec(fmt.Sprintf("DELETE FROM %s", table))
	}
	
	// 插入课程
	fmt.Println("  插入课程数据...")
	courses := []struct {
		title, description, cover string
		lessonNum, limitCount     int
		credit                    float64
	}{
		{
			"网络安全基础", 
			"学习网络安全的基本概念和原理，包括加密、认证、访问控制等",
			"https://example.com/covers/security-basic.jpg",
			15, 100, 3.0,
		},
		{
			"渗透测试实战", 
			"通过实战演练学习渗透测试的方法和工具，掌握漏洞挖掘技巧",
			"https://example.com/covers/penetration-test.jpg",
			20, 80, 4.0,
		},
		{
			"密码学原理与应用", 
			"深入理解密码学原理，学习对称加密、非对称加密、哈希算法等",
			"https://example.com/covers/cryptography.jpg",
			18, 120, 3.5,
		},
		{
			"Web安全攻防", 
			"学习Web安全常见漏洞和防护措施，包括SQL注入、XSS、CSRF等",
			"https://example.com/covers/web-security.jpg",
			16, 90, 3.0,
		},
	}
	
	for i, course := range courses {
		_, err := db.Exec(
			"INSERT INTO courses (title, description, cover, lesson_num, credit, limit_count) VALUES (?, ?, ?, ?, ?, ?)",
			course.title, course.description, course.cover, course.lessonNum, course.credit, course.limitCount,
		)
		if err != nil {
			log.Printf("插入课程 %d 失败: %v", i+1, err)
		}
	}
	
	// 插入教师
	fmt.Println("  插入教师数据...")
	teachers := []struct {
		name, career, intro string
	}{
		{
			"张教授",
			"网络安全专家",
			"从事网络安全研究20年，发表多篇SCI论文，有丰富的教学和实践经验",
		},
		{
			"李博士",
			"密码学研究员",
			"专注于密码学算法研究，参与多个国家级安全项目",
		},
	}
	
	for i, teacher := range teachers {
		result, err := db.Exec(
			"INSERT INTO teachers (name, career, intro) VALUES (?, ?, ?)",
			teacher.name, teacher.career, teacher.intro,
		)
		if err != nil {
			log.Printf("插入教师 %d 失败: %v", i+1, err)
		} else {
			// 关联教师和课程
			teacherID, _ := result.LastInsertId()
			_, err = db.Exec(
				"INSERT INTO teacher_courses (teacher_id, course_id) VALUES (?, ?)",
				teacherID, i+1, // 第一个教师教第一门课，第二个教师教第二门课
			)
			if err != nil {
				log.Printf("关联教师课程失败: %v", err)
			}
		}
	}
	
	// 插入视频
	fmt.Println("  插入视频数据...")
	videos := []struct {
		url, description string
		duration         int
	}{
		{"http://localhost:3000/api/videoing/security1.mp4", "网络安全概述", 3600},
		{"http://localhost:3000/api/videoing/security2.mp4", "加密技术基础", 4200},
		{"http://localhost:3000/api/videoing/penetration1.mp4", "渗透测试入门", 3800},
		{"http://localhost:3000/api/videoing/penetration2.mp4", "漏洞扫描工具", 4500},
	}
	
	for i, video := range videos {
		_, err := db.Exec(
			"INSERT INTO videos (url, description, duration) VALUES (?, ?, ?)",
			video.url, video.description, video.duration,
		)
		if err != nil {
			log.Printf("插入视频 %d 失败: %v", i+1, err)
		}
	}
	
	// 插入章节和课时
	fmt.Println("  插入章节数据...")
	
	// 第一门课的章节
	db.Exec("INSERT INTO chapters (course_id, title, state) VALUES (1, '第一章：网络安全基础', 1)")
	db.Exec("INSERT INTO chapters (course_id, title, state) VALUES (1, '第二章：加密技术', 0)")
	
	// 第二门课的章节
	db.Exec("INSERT INTO chapters (course_id, title, state) VALUES (2, '第一章：渗透测试概述', 1)")
	db.Exec("INSERT INTO chapters (course_id, title, state) VALUES (2, '第二章：信息收集', 0)")
	
	// 第一门课的课时
	db.Exec("INSERT INTO chapter_children (chapter_id, title, video_id) VALUES (1, '1.1 网络安全概念', 1)")
	db.Exec("INSERT INTO chapter_children (chapter_id, title, video_id) VALUES (1, '1.2 安全威胁分析', 2)")
	db.Exec("INSERT INTO chapter_children (chapter_id, title, video_id) VALUES (2, '2.1 对称加密', 2)")
	
	// 第二门课的课时
	db.Exec("INSERT INTO chapter_children (chapter_id, title, video_id) VALUES (3, '1.1 渗透测试流程', 3)")
	db.Exec("INSERT INTO chapter_children (chapter_id, title, video_id) VALUES (3, '1.2 法律与道德', 3)")
	db.Exec("INSERT INTO chapter_children (chapter_id, title, video_id) VALUES (4, '2.1 信息收集方法', 4)")
	
	fmt.Println("✓ 测试数据插入完成")
}

func verifyData(db *sql.DB) {
	// 统计各表数据量
	tables := []string{"courses", "teachers", "videos", "chapters", "chapter_children"}
	
	fmt.Println("┌──────────────┬────────┐")
	fmt.Println("│ 表名         │ 记录数 │")
	fmt.Println("├──────────────┼────────┤")
	
	for _, table := range tables {
		var count int
		err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			fmt.Printf("│ %-12s │ 查询失败 │\n", table)
		} else {
			fmt.Printf("│ %-12s │ %6d │\n", table, count)
		}
	}
	fmt.Println("└──────────────┴────────┘")
	
	// 显示课程列表
	fmt.Println("\n课程列表:")
	rows, err := db.Query("SELECT id, title, lesson_num, credit FROM courses ORDER BY id")
	if err != nil {
		log.Printf("查询课程失败: %v", err)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var id, lessonNum int
		var title string
		var credit float64
		err := rows.Scan(&id, &title, &lessonNum, &credit)
		if err != nil {
			log.Printf("读取课程数据失败: %v", err)
			continue
		}
		fmt.Printf("  ID:%d 《%s》 课时:%d 学分:%.1f\n", id, title, lessonNum, credit)
	}
}