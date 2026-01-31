package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("=== 更新数据库密码为bcrypt哈希 ===")
	
	// 数据库连接
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
	
	// 生成bcrypt哈希
	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成bcrypt哈希失败:", err)
	}
	
	hashedStr := string(hashedPassword)
	fmt.Printf("✓ 密码 '%s' 的bcrypt哈希已生成\n", password)
	fmt.Printf("  哈希值: %s\n", hashedStr)
	
	// 更新数据库中的密码
	fmt.Println("\n更新用户密码...")
	
	updateSQL := "UPDATE students SET password = ? WHERE stuId IN ('20230001', '20230002')"
	result, err := db.Exec(updateSQL, hashedStr)
	if err != nil {
		log.Fatal("更新密码失败:", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("✓ 更新了 %d 个用户的密码\n", rowsAffected)
	
	// 验证更新
	fmt.Println("\n验证更新结果:")
	query := "SELECT stuId, LENGTH(password) as pwd_len FROM students WHERE stuId IN ('20230001', '20230002')"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	defer rows.Close()
	
	fmt.Println("学号\t\t密码长度")
	fmt.Println("----\t\t--------")
	for rows.Next() {
		var stuID string
		var pwdLen int
		err := rows.Scan(&stuID, &pwdLen)
		if err != nil {
			log.Fatal("读取失败:", err)
		}
		fmt.Printf("%s\t\t%d\n", stuID, pwdLen)
		
		// bcrypt哈希通常是60个字符
		if pwdLen >= 50 && pwdLen <= 60 {
			fmt.Printf("  ✅ 密码已正确存储为bcrypt哈希\n")
		} else {
			fmt.Printf("  ⚠️  密码长度异常，可能不是bcrypt哈希\n")
		}
	}
	
	fmt.Println("\n✅ 密码更新完成！")
	fmt.Println("现在可以使用bcrypt验证登录了。")
}