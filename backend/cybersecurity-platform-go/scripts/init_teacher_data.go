package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=== 初始化教师数据 ===")

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

	// 检查teachers表
	fmt.Println("\n1. 检查teachers表...")
	checkTeachersTable(db)

	// 插入测试教师数据
	fmt.Println("\n2. 插入测试教师数据...")
	insertTeacherTestData(db)

	// 验证数据
	fmt.Println("\n3. 验证数据...")
	verifyTeacherData(db)

	fmt.Println("\n✅ 教师数据初始化完成！")
}

func checkTeachersTable(db *sql.DB) {
	// 检查teachers表是否存在
	var tableExists string
	err := db.QueryRow("SHOW TABLES LIKE 'teachers'").Scan(&tableExists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("检查teachers表失败: %v", err)
	}

	if tableExists == "" {
		fmt.Println("⚠️  teachers表不存在")
		fmt.Println("   请先运行课程数据初始化脚本")
	} else {
		fmt.Println("✓ teachers表已存在")
		
		// 检查表结构
		rows, err := db.Query("DESCRIBE teachers")
		if err != nil {
			log.Printf("查看表结构失败: %v", err)
		} else {
			fmt.Println("  teachers表结构:")
			for rows.Next() {
				var field, fieldType, null, key, defaultValue, extra string
				err := rows.Scan(&field, &fieldType, &null, &key, &defaultValue, &extra)
				if err != nil {
					log.Printf("读取表结构失败: %v", err)
					continue
				}
				fmt.Printf("    字段: %-15s 类型: %-20s\n", field, fieldType)
			}
			rows.Close()
		}
	}
}

func insertTeacherTestData(db *sql.DB) {
	// 检查是否已有教师数据
	var teacherCount int
	err := db.QueryRow("SELECT COUNT(*) FROM teachers").Scan(&teacherCount)
	if err != nil {
		log.Printf("查询教师数量失败: %v", err)
		return
	}

	if teacherCount == 0 {
		fmt.Println("  没有教师数据，插入测试教师...")
		
		// 插入教师数据
		teachers := []struct {
			name   string
			career string
			intro  string
		}{
			{
				"张教授",
				"网络安全专家，10年安全研究经验",
				"专注于网络安全教育，曾参与多个国家级安全项目",
			},
			{
				"李老师",
				"渗透测试工程师，OWASP贡献者",
				"具有丰富的实战经验，擅长Web安全测试",
			},
			{
				"王博士",
				"密码学研究员，加密算法专家",
				"在国内外期刊发表多篇密码学论文",
			},
			{
				"赵教授",
				"数据安全专家，数据库安全顾问",
				"专注于数据库安全和数据保护技术",
			},
			{
				"孙老师",
				"移动安全工程师，Android安全专家",
				"擅长移动应用安全测试和防护",
			},
		}

		for i, teacher := range teachers {
			_, err := db.Exec(
				"INSERT INTO teachers (name, career, intro) VALUES (?, ?, ?)",
				teacher.name, teacher.career, teacher.intro,
			)
			if err != nil {
				log.Printf("插入教师 %d 失败: %v", i+1, err)
			} else {
				fmt.Printf("✓ 插入教师: %s\n", teacher.name)
			}
		}
	} else {
		fmt.Printf("✓ 已有 %d 位教师数据\n", teacherCount)
	}

	// 检查教师-课程关联
	fmt.Println("\n  检查教师-课程关联...")
	
	// 获取所有教师
	rows, err := db.Query("SELECT id, name FROM teachers")
	if err != nil {
		log.Printf("查询教师列表失败: %v", err)
		return
	}
	defer rows.Close()

	teachers := make(map[int]string)
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Printf("读取教师数据失败: %v", err)
			continue
		}
		teachers[id] = name
	}

	// 获取所有课程
	courseRows, err := db.Query("SELECT id, title FROM courses")
	if err != nil {
		log.Printf("查询课程列表失败: %v", err)
		return
	}
	defer courseRows.Close()

	courses := []struct {
		id    int
		title string
	}{}
	for courseRows.Next() {
		var id int
		var title string
		err := courseRows.Scan(&id, &title)
		if err != nil {
			log.Printf("读取课程数据失败: %v", err)
			continue
		}
		courses = append(courses, struct {
			id    int
			title string
		}{id, title})
	}

	// 为每个课程分配一个教师（循环分配）
	if len(teachers) > 0 && len(courses) > 0 {
		teacherIDs := make([]int, 0, len(teachers))
		for id := range teachers {
			teacherIDs = append(teacherIDs, id)
		}

		for i, course := range courses {
			teacherIndex := i % len(teacherIDs)
			teacherID := teacherIDs[teacherIndex]

			// 检查关联是否已存在
			var exists bool
			err := db.QueryRow(
				"SELECT EXISTS(SELECT 1 FROM teacher_courses WHERE teacher_id = ? AND course_id = ?)",
				teacherID, course.id,
			).Scan(&exists)

			if err != nil {
				log.Printf("检查关联存在失败: %v", err)
				continue
			}

			if !exists {
				_, err := db.Exec(
					"INSERT INTO teacher_courses (teacher_id, course_id) VALUES (?, ?)",
					teacherID, course.id,
				)
				if err != nil {
					log.Printf("关联教师课程失败: %v", err)
				} else {
					fmt.Printf("✓ 关联课程[%s]与教师[%s]\n", course.title, teachers[teacherID])
				}
			}
		}
	}
}

func verifyTeacherData(db *sql.DB) {
	// 统计教师数量
	var teacherCount int
	err := db.QueryRow("SELECT COUNT(*) FROM teachers").Scan(&teacherCount)
	if err != nil {
		log.Printf("统计教师数量失败: %v", err)
		return
	}

	fmt.Printf("教师总数: %d\n", teacherCount)

	// 显示教师列表
	fmt.Println("\n教师列表:")
	rows, err := db.Query(`
		SELECT t.id, t.name, t.career, 
		       COUNT(tc.course_id) as course_count
		FROM teachers t
		LEFT JOIN teacher_courses tc ON t.id = tc.teacher_id
		GROUP BY t.id
		ORDER BY t.id
	`)
	if err != nil {
		log.Printf("查询教师列表失败: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, courseCount int
		var name, career string
		err := rows.Scan(&id, &name, &career, &courseCount)
		if err != nil {
			log.Printf("读取教师数据失败: %v", err)
			continue
		}
		fmt.Printf("  ID:%d %s - %s (课程数: %d)\n", id, name, career, courseCount)
	}

	// 显示教师-课程关联详情
	fmt.Println("\n教师-课程关联详情:")
	detailRows, err := db.Query(`
		SELECT t.name, c.title
		FROM teacher_courses tc
		JOIN teachers t ON tc.teacher_id = t.id
		JOIN courses c ON tc.course_id = c.id
		ORDER BY t.name, c.title
	`)
	if err != nil {
		log.Printf("查询关联详情失败: %v", err)
		return
	}
	defer detailRows.Close()

	for detailRows.Next() {
		var teacherName, courseTitle string
		err := detailRows.Scan(&teacherName, &courseTitle)
		if err != nil {
			log.Printf("读取关联详情失败: %v", err)
			continue
		}
		fmt.Printf("  %s → 《%s》\n", teacherName, courseTitle)
	}
}