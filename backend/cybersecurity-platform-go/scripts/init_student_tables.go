package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=== 初始化学生选课表 ===")

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

	// 创建学生选课表
	fmt.Println("\n1. 创建学生选课表...")
	createStudentTables(db)

	// 插入测试数据
	fmt.Println("\n2. 插入测试数据...")
	insertStudentTestData(db)

	// 验证数据
	fmt.Println("\n3. 验证数据...")
	verifyStudentData(db)

	fmt.Println("\n✅ 学生选课表初始化完成！")
}

func createStudentTables(db *sql.DB) {
	// 创建学生选课表
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS student_courses (
			id INT PRIMARY KEY AUTO_INCREMENT,
			stuId VARCHAR(50) NOT NULL,
			course_id INT NOT NULL,
			joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
			UNIQUE KEY unique_student_course (stuId, course_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("创建student_courses表失败: %v", err)
	} else {
		fmt.Println("✓ student_courses表创建成功")
	}

	// 检查students表是否存在（从登录注册模块）
	checkStudentsTableSQL := "SHOW TABLES LIKE 'students'"
	var tableExists string
	err = db.QueryRow(checkStudentsTableSQL).Scan(&tableExists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("检查students表失败: %v", err)
	}

	if tableExists == "" {
		fmt.Println("⚠️  students表不存在，需要先运行用户注册/登录初始化")
	} else {
		fmt.Println("✓ students表已存在")
	}

	// 检查userdetail表是否存在
	checkUserDetailSQL := "SHOW TABLES LIKE 'userdetail'"
	err = db.QueryRow(checkUserDetailSQL).Scan(&tableExists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("检查userdetail表失败: %v", err)
	}

	if tableExists == "" {
		fmt.Println("⚠️  userdetail表不存在，需要先运行用户注册/登录初始化")
	} else {
		fmt.Println("✓ userdetail表已存在")
	}
}

func insertStudentTestData(db *sql.DB) {
	// 检查是否有测试学生
	fmt.Println("  检查测试学生数据...")

	// 注意：这里假设students和userdetail表已经通过登录/注册模块创建
	// 我们只检查是否存在，不重复插入

	var studentCount int
	err := db.QueryRow("SELECT COUNT(*) FROM students WHERE stuId IN ('20230001', '20230002')").Scan(&studentCount)
	if err != nil {
		log.Printf("查询测试学生失败: %v", err)
	}

	if studentCount == 0 {
		fmt.Println("⚠️  测试学生不存在，请先运行用户注册/登录初始化")
		fmt.Println("   或使用以下学号进行测试：20230001, 20230002")
	} else {
		fmt.Printf("✓ 找到 %d 个测试学生\n", studentCount)
	}

	// 插入测试选课记录（如果不存在）
	fmt.Println("  插入测试选课记录...")

	// 获取现有的课程
	rows, err := db.Query("SELECT id, title FROM courses ORDER BY id LIMIT 3")
	if err != nil {
		log.Printf("查询课程失败: %v", err)
		return
	}
	defer rows.Close()

	courseIDs := []int{}
	courseTitles := []string{}
	for rows.Next() {
		var id int
		var title string
		err := rows.Scan(&id, &title)
		if err != nil {
			log.Printf("读取课程数据失败: %v", err)
			continue
		}
		courseIDs = append(courseIDs, id)
		courseTitles = append(courseTitles, title)
	}

	if len(courseIDs) == 0 {
		fmt.Println("⚠️  没有找到课程，请先运行课程数据初始化")
		return
	}

	// 为测试学生添加选课记录
	testEnrollments := []struct {
		stuId    string
		courseID int
	}{
		{"20230001", courseIDs[0]},
		{"20230001", courseIDs[1]},
		{"20230002", courseIDs[0]},
	}

	enrollmentCount := 0
	for _, enrollment := range testEnrollments {
		// 检查是否已存在
		var exists bool
		err := db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM student_courses WHERE stuId = ? AND course_id = ?)",
			enrollment.stuId, enrollment.courseID,
		).Scan(&exists)

		if err != nil {
			log.Printf("检查选课记录失败: %v", err)
			continue
		}

		if !exists {
			_, err := db.Exec(
				"INSERT INTO student_courses (stuId, course_id) VALUES (?, ?)",
				enrollment.stuId, enrollment.courseID,
			)
			if err != nil {
				log.Printf("插入选课记录失败: %v", err)
			} else {
				enrollmentCount++
				fmt.Printf("✓ 学生 %s 加入课程 ID:%d\n", enrollment.stuId, enrollment.courseID)
			}
		} else {
			fmt.Printf("⏭️  选课记录已存在: 学生 %s - 课程 ID:%d\n", enrollment.stuId, enrollment.courseID)
		}
	}

	if enrollmentCount > 0 {
		fmt.Printf("✓ 新增 %d 条选课记录\n", enrollmentCount)
	} else {
		fmt.Println("⏭️  所有选课记录已存在")
	}
}

func verifyStudentData(db *sql.DB) {
	// 统计选课记录
	var enrollmentCount int
	err := db.QueryRow("SELECT COUNT(*) FROM student_courses").Scan(&enrollmentCount)
	if err != nil {
		log.Printf("统计选课记录失败: %v", err)
		return
	}

	fmt.Printf("选课记录总数: %d\n", enrollmentCount)

	// 显示选课详情
	fmt.Println("\n选课详情:")
	rows, err := db.Query(`
		SELECT sc.stuId, u.nickName, c.title, sc.joined_at 
		FROM student_courses sc
		JOIN courses c ON sc.course_id = c.id
		LEFT JOIN userdetail u ON sc.stuId = u.stuId
		ORDER BY sc.joined_at DESC
	`)
	if err != nil {
		log.Printf("查询选课详情失败: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var stuID, nickName, courseTitle, joinedAt string
		err := rows.Scan(&stuID, &nickName, &courseTitle, &joinedAt)
		if err != nil {
			log.Printf("读取选课详情失败: %v", err)
			continue
		}
		if nickName == "" {
			nickName = "未知用户"
		}
		fmt.Printf("  %s (%s) - 《%s》 - 加入时间: %s\n", nickName, stuID, courseTitle, joinedAt)
	}
}