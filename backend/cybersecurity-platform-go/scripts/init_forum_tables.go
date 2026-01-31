package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=== 初始化论坛数据表 ===")

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

	// 创建论坛相关表
	fmt.Println("\n1. 创建论坛相关表...")
	createForumTables(db)

	// 插入测试数据
	fmt.Println("\n2. 插入测试数据...")
	insertForumTestData(db)

	// 验证数据
	fmt.Println("\n3. 验证数据...")
	verifyForumData(db)

	fmt.Println("\n✅ 论坛数据表初始化完成！")
}

func createForumTables(db *sql.DB) {
	tables := []struct {
		name string
		sql  string
	}{
		{
			name: "forum_categories",
			sql: `
				CREATE TABLE IF NOT EXISTS forum_categories (
					id INT PRIMARY KEY AUTO_INCREMENT,
					name VARCHAR(50) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
			`,
		},
		{
			name: "forum_tags",
			sql: `
				CREATE TABLE IF NOT EXISTS forum_tags (
					id INT PRIMARY KEY AUTO_INCREMENT,
					name VARCHAR(50) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
			`,
		},
		{
			name: "forum_articles",
			sql: `
				CREATE TABLE IF NOT EXISTS forum_articles (
					id INT PRIMARY KEY AUTO_INCREMENT,
					title VARCHAR(200) NOT NULL,
					stuId VARCHAR(50) NOT NULL,
					stuName VARCHAR(50) NOT NULL,
					stuHead VARCHAR(500),
					cateId INT DEFAULT 0,
					isTop TINYINT(1) DEFAULT 0,
					isEss TINYINT(1) DEFAULT 0,
					viewCount INT DEFAULT 0,
					createTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					FOREIGN KEY (cateId) REFERENCES forum_categories(id)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
			`,
		},
		{
			name: "article_tags",
			sql: `
				CREATE TABLE IF NOT EXISTS article_tags (
					id INT PRIMARY KEY AUTO_INCREMENT,
					article_id INT NOT NULL,
					tag_id INT NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (article_id) REFERENCES forum_articles(id) ON DELETE CASCADE,
					FOREIGN KEY (tag_id) REFERENCES forum_tags(id),
					UNIQUE KEY unique_article_tag (article_id, tag_id)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
			`,
		},
		{
			name: "forum_comments",
			sql: `
				CREATE TABLE IF NOT EXISTS forum_comments (
					id INT PRIMARY KEY AUTO_INCREMENT,
					article_id INT NOT NULL,
					content TEXT NOT NULL,
					author_id VARCHAR(50) NOT NULL,
					author_name VARCHAR(50) NOT NULL,
					author_head VARCHAR(500),
					parent_id INT DEFAULT NULL,
					status TINYINT(1) DEFAULT 1,
					like_count INT DEFAULT 0,
					create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (article_id) REFERENCES forum_articles(id) ON DELETE CASCADE
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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

func insertForumTestData(db *sql.DB) {
	// 插入分类
	fmt.Println("  插入分类数据...")
	categories := []string{"网络安全", "漏洞分析", "安全工具", "密码学", "法律法规"}
	
	for _, category := range categories {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM forum_categories WHERE name = ?)", category).Scan(&exists)
		if err != nil {
			log.Printf("检查分类存在失败: %v", err)
			continue
		}
		
		if !exists {
			_, err := db.Exec("INSERT INTO forum_categories (name) VALUES (?)", category)
			if err != nil {
				log.Printf("插入分类 %s 失败: %v", category, err)
			} else {
				fmt.Printf("✓ 插入分类: %s\n", category)
			}
		}
	}

	// 插入标签
	fmt.Println("  插入标签数据...")
	tags := []string{"渗透测试", "Web安全", "移动安全", "云安全", "数据加密", "身份认证", "防火墙", "入侵检测"}
	
	for _, tag := range tags {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM forum_tags WHERE name = ?)", tag).Scan(&exists)
		if err != nil {
			log.Printf("检查标签存在失败: %v", err)
			continue
		}
		
		if !exists {
			_, err := db.Exec("INSERT INTO forum_tags (name) VALUES (?)", tag)
			if err != nil {
				log.Printf("插入标签 %s 失败: %v", tag, err)
			} else {
				fmt.Printf("✓ 插入标签: %s\n", tag)
			}
		}
	}
	
	fmt.Println("✓ 测试数据插入完成")
	fmt.Println("  注：文章和评论数据将在后续阶段插入")
}

func verifyForumData(db *sql.DB) {
	// 统计各表数据量
	tables := []string{"forum_categories", "forum_tags", "forum_articles", "article_tags", "forum_comments"}
	
	fmt.Println("┌────────────────────┬────────┐")
	fmt.Println("│ 表名               │ 记录数 │")
	fmt.Println("├────────────────────┼────────┤")
	
	for _, table := range tables {
		var count int
		err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			fmt.Printf("│ %-18s │ 查询失败 │\n", table)
		} else {
			fmt.Printf("│ %-18s │ %6d │\n", table, count)
		}
	}
	fmt.Println("└────────────────────┴────────┘")
	
	// 显示分类列表
	fmt.Println("\n分类列表:")
	rows, err := db.Query("SELECT id, name FROM forum_categories ORDER BY id")
	if err != nil {
		log.Printf("查询分类失败: %v", err)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Printf("读取分类数据失败: %v", err)
			continue
		}
		fmt.Printf("  ID:%d %s\n", id, name)
	}
	
	// 显示标签列表
	fmt.Println("\n标签列表:")
	tagRows, err := db.Query("SELECT id, name FROM forum_tags ORDER BY id LIMIT 10")
	if err != nil {
		log.Printf("查询标签失败: %v", err)
		return
	}
	defer tagRows.Close()
	
	for tagRows.Next() {
		var id int
		var name string
		err := tagRows.Scan(&id, &name)
		if err != nil {
			log.Printf("读取标签数据失败: %v", err)
			continue
		}
		fmt.Printf("  ID:%d %s\n", id, name)
	}
}