package main

import (
	"database/sql"
	"fmt"
	"log"

	"cybersecurity-platform-go/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=== 初始化课程数据 ===")
	
	// 加载配置
	cfg := config.LoadConfig()
	
	// 连接到数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()
	
	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}
	
	fmt.Println("✓ 数据库连接成功")
	
	// 1. 检查并创建courses表
	fmt.Println("\n1. 检查courses表...")
	createCoursesTable(db)
	
	// 2. 检查并创建teachers表
	fmt.Println("\n2. 检查teachers表...")
	createTeachersTable(db)
	
	// 3. 检查并创建teacher_courses表
	fmt.Println("\n3. 检查teacher_courses表...")
	createTeacherCoursesTable(db)
	
	// 4. 插入测试数据
	fmt.Println("\n4. 插入测试数据...")
	insertTestData(db)
	
	fmt.Println("\n✅ 课程数据初始化完成！")
}

func createCoursesTable(db *sql.DB) {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS courses (
			id INT PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			cover VARCHAR(500),
			lesson_num INT DEFAULT 0,
			credit INT DEFAULT 0,
			limit_count INT DEFAULT 100,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("创建courses表失败: %v", err)
	} else {
		fmt.Println("✓ courses表检查/创建完成")
	}
}

func createTeachersTable(db *sql.DB) {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS teachers (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(100) NOT NULL,
			career VARCHAR(200),
			intro TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("创建teachers表失败: %v", err)
	} else {
		fmt.Println("✓ teachers表检查/创建完成")
	}
}

func createTeacherCoursesTable(db *sql.DB) {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS teacher_courses (
			id INT PRIMARY KEY AUTO_INCREMENT,
			teacher_id INT NOT NULL,
			course_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (teacher_id) REFERENCES teachers(id),
			FOREIGN KEY (course_id) REFERENCES courses(id),
			UNIQUE KEY unique_teacher_course (teacher_id, course_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("创建teacher_courses表失败: %v", err)
	} else {
		fmt.Println("✓ teacher_courses表检查/创建完成")
	}
}

func insertTestData(db *sql.DB) {
	// 插入课程数据
	courses := []struct {
		title       string
		description string
		cover       string
		lessonNum   int
		credit      int
		limitCount  int
	}{
		{
			"网络安全基础入门",
			"学习网络安全的基本概念、原理和技术，为后续学习打下坚实基础",
			"/img/course/security-basic.jpg",
			12,
			3,
			50,
		},
		{
			"渗透测试实战",
			"掌握渗透测试的完整流程，包括信息收集、漏洞扫描、利用和报告编写",
			"/img/course/penetration-test.jpg",
			15,
			4,
			40,
		},
		{
			"数据加密与解密",
			"学习对称加密、非对称加密算法原理，以及各种加密技术的实际应用",
			"/img/course/data-encryption.jpg",
			10,
			3,
			60,
		},
		{
			"Web安全攻防",
			"深入理解Web应用常见漏洞，学习SQL注入、XSS、CSRF等攻击的防御方法",
			"/img/course/web-security.jpg",
			14,
			4,
			45,
		},
		{
			"移动应用安全",
			"学习Android和iOS应用的安全测试方法，掌握移动端安全防护技术",
			"/img/course/mobile-security.jpg",
			11,
			3,
			55,
		},
	}
	
	for _, course := range courses {
		// 检查课程是否已存在
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM courses WHERE title = ?", course.title).Scan(&count)
		if err != nil {
			log.Printf("检查课程存在失败: %v", err)
			continue
		}
		
		if count == 0 {
			_, err := db.Exec(`
				INSERT INTO courses (title, description, cover, lesson_num, credit, limit_count) 
				VALUES (?, ?, ?, ?, ?, ?)
			`, course.title, course.description, course.cover, course.lessonNum, course.credit, course.limitCount)
			
			if err != nil {
				log.Printf("插入课程失败: %v", err)
			} else {
				fmt.Printf("✓ 插入课程: %s\n", course.title)
			}
		} else {
			fmt.Printf("⏭️  课程已存在: %s\n", course.title)
		}
	}
	
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
	}
	
	teacherIDs := []int{}
	for _, teacher := range teachers {
		// 检查教师是否已存在
		var id int
		err := db.QueryRow("SELECT id FROM teachers WHERE name = ?", teacher.name).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("检查教师存在失败: %v", err)
			continue
		}
		
		if err == sql.ErrNoRows {
			result, err := db.Exec(`
				INSERT INTO teachers (name, career, intro) 
				VALUES (?, ?, ?)
			`, teacher.name, teacher.career, teacher.intro)
			
			if err != nil {
				log.Printf("插入教师失败: %v", err)
				continue
			}
			
			lastID, _ := result.LastInsertId()
			teacherIDs = append(teacherIDs, int(lastID))
			fmt.Printf("✓ 插入教师: %s\n", teacher.name)
		} else {
			teacherIDs = append(teacherIDs, id)
			fmt.Printf("⏭️  教师已存在: %s (ID: %d)\n", teacher.name, id)
		}
	}
	
	// 关联教师和课程
	if len(teacherIDs) > 0 {
		// 获取所有课程ID
		rows, err := db.Query("SELECT id, title FROM courses ORDER BY id")
		if err != nil {
			log.Printf("获取课程列表失败: %v", err)
			return
		}
		defer rows.Close()
		
		var courseIDs []int
		var courseTitles []string
		for rows.Next() {
			var id int
			var title string
			err := rows.Scan(&id, &title)
			if err != nil {
				log.Printf("扫描课程数据失败: %v", err)
				continue
			}
			courseIDs = append(courseIDs, id)
			courseTitles = append(courseTitles, title)
		}
		
		// 为每个课程分配一个教师
		for i, courseID := range courseIDs {
			teacherIndex := i % len(teacherIDs)
			teacherID := teacherIDs[teacherIndex]
			
			// 检查关联是否已存在
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM teacher_courses WHERE teacher_id = ? AND course_id = ?", 
				teacherID, courseID).Scan(&count)
			
			if err != nil {
				log.Printf("检查关联存在失败: %v", err)
				continue
			}
			
			if count == 0 {
				_, err := db.Exec(`
					INSERT INTO teacher_courses (teacher_id, course_id) 
					VALUES (?, ?)
				`, teacherID, courseID)
				
				if err != nil {
					log.Printf("关联教师课程失败: %v", err)
				} else {
					fmt.Printf("✓ 关联课程[%s]与教师[%d]\n", courseTitles[i], teacherID)
				}
			}
		}
	}
	
	fmt.Println("\n✅ 测试数据插入完成")
	fmt.Println("课程总数:", len(courses))
	fmt.Println("教师总数:", len(teachers))
}