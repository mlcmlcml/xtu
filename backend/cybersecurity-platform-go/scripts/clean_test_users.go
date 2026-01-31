package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=== 清理测试用户数据 ===")
	
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
	
	// 删除测试注册用户（保留20230001和20230002用于登录测试）
	deleteSQL := `
		DELETE u FROM userdetail u
		JOIN students s ON u.stuId = s.stuId
		WHERE s.stuId NOT IN ('20230001', '20230002')
	`
	
	result, err := db.Exec(deleteSQL)
	if err != nil {
		log.Fatal("删除userdetail数据失败:", err)
	}
	
	rowsAffected1, _ := result.RowsAffected()
	fmt.Printf("✓ 删除了 %d 条userdetail记录\n", rowsAffected1)
	
	// 删除students表中的对应记录
	deleteStudentsSQL := "DELETE FROM students WHERE stuId NOT IN ('20230001', '20230002')"
	result, err = db.Exec(deleteStudentsSQL)
	if err != nil {
		log.Fatal("删除students数据失败:", err)
	}
	
	rowsAffected2, _ := result.RowsAffected()
	fmt.Printf("✓ 删除了 %d 条students记录\n", rowsAffected2)
	
	// 显示剩余用户
	fmt.Println("\n剩余用户:")
	query := `
		SELECT s.stuId, s.email, u.nickName
		FROM students s
		LEFT JOIN userdetail u ON s.stuId = u.stuId
		ORDER BY s.stuId
	`
	
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	defer rows.Close()
	
	count := 0
	for rows.Next() {
		var stuID, email, nickName string
		err := rows.Scan(&stuID, &email, &nickName)
		if err != nil {
			log.Fatal("读取失败:", err)
		}
		fmt.Printf("  学号: %s, 邮箱: %s, 昵称: %s\n", stuID, email, nickName)
		count++
	}
	
	fmt.Printf("\n✅ 清理完成！共剩余 %d 个用户\n", count)
}