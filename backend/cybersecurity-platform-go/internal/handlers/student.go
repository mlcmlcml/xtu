// internal/handlers/student.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"cybersecurity-platform-go/internal/database"
)

// StudentJoinRequest 学生加入课程请求
type StudentJoinRequest struct {
	CourseID int    `json:"courseId"`
	StuID    string `json:"stuId"`
}

// StudentJoinResponse 学生加入课程响应
type StudentJoinResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MyCourse 我的课程信息
type MyCourse struct {
	ID          int    `json:"id"`
	ClassName   string `json:"classname"`
	TeacherName string `json:"teachername"`
	Count1      int    `json:"count1"`  // 课时数
	Count2      int    `json:"count2"`  // 限制人数
	Cover       string `json:"cover"`
	Career      string `json:"career"`
}

// MyCoursesResponse 我的课程响应
type MyCoursesResponse struct {
	Code int `json:"code"`
	Data struct {
		Courses []MyCourse `json:"courses"`
		Total   int        `json:"total"`
	} `json:"data"`
}

// CheckEnrollmentResponse 检查选课状态响应
type CheckEnrollmentResponse struct {
	Code int `json:"code"`
	Data struct {
		IsEnrolled bool `json:"isEnrolled"`
	} `json:"data"`
}

// RegisterStudentRoutes 注册学生路由
func RegisterStudentRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// 加入课程
	mux.HandleFunc("POST /api/student/joinCourse", joinCourseHandler)

	// 获取我的课程
	mux.HandleFunc("GET /api/student/myCourses", myCoursesHandler)

	// 检查选课状态
	mux.HandleFunc("GET /api/student/checkEnrollment", checkEnrollmentHandler)

	return mux
}

// joinCourseHandler 加入课程处理器
func joinCourseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 只处理POST请求
	if r.Method != http.MethodPost {
		sendStudentError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		return
	}

	// 解析请求体
	var req StudentJoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendStudentError(w, http.StatusBadRequest, 40000, "参数解析失败")
		return
	}

	// 调用加入课程逻辑
	joinCourse(w, req)
}

// myCoursesHandler 获取我的课程处理器
func myCoursesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 只处理GET请求
	if r.Method != http.MethodGet {
		sendStudentError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		return
	}

	// 解析查询参数
	stuID := r.URL.Query().Get("stuId")
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	if stuID == "" {
		sendStudentError(w, http.StatusBadRequest, 40000, "缺少学生ID参数")
		return
	}

	// 设置默认值
	page := 1
	size := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			size = s
		}
	}

	// 调用获取我的课程逻辑
	myCourses(w, stuID, page, size)
}

// checkEnrollmentHandler 检查选课状态处理器
func checkEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 只处理GET请求
	if r.Method != http.MethodGet {
		sendStudentError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		return
	}

	// 解析查询参数
	courseIDStr := r.URL.Query().Get("courseId")
	stuID := r.URL.Query().Get("stuId")

	if courseIDStr == "" || stuID == "" {
		sendStudentError(w, http.StatusBadRequest, 40000, "参数不完整")
		return
	}

	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil || courseID <= 0 {
		sendStudentError(w, http.StatusBadRequest, 40000, "无效的课程ID")
		return
	}

	// 调用检查选课状态逻辑
	checkEnrollment(w, courseID, stuID)
}

// sendStudentError 发送学生相关错误响应
func sendStudentError(w http.ResponseWriter, httpStatus, code int, message string) {
	w.WriteHeader(httpStatus)
	response := map[string]interface{}{
		"code":    code,
		"message": message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码错误响应失败: %v", err)
	}
}

// joinCourse 加入课程的核心逻辑
func joinCourse(w http.ResponseWriter, req StudentJoinRequest) {
	// 检查参数是否完整
	if req.CourseID == 0 || req.StuID == "" {
		sendStudentError(w, http.StatusBadRequest, 40000, "参数不完整")
		return
	}

	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 检查课程是否存在
	var courseExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM courses WHERE id = ?)", req.CourseID).Scan(&courseExists)
	if err != nil {
		log.Printf("检查课程存在失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	if !courseExists {
		sendStudentError(w, http.StatusNotFound, 40400, "课程不存在")
		return
	}

	// 检查学生是否存在
	var studentExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM students WHERE stuId = ?)", req.StuID).Scan(&studentExists)
	if err != nil {
		log.Printf("检查学生存在失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	if !studentExists {
		sendStudentError(w, http.StatusNotFound, 40401, "学生不存在")
		return
	}

	// 检查是否已加入
	var alreadyJoined bool
	err = db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM student_courses WHERE stuId = ? AND course_id = ?)",
		req.StuID, req.CourseID,
	).Scan(&alreadyJoined)

	if err != nil {
		log.Printf("检查选课状态失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	if alreadyJoined {
		response := StudentJoinResponse{
			Code:    40001,
			Message: "已加入该课程",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查课程人数是否已满
	var limitCount, currentCount int
	err = db.QueryRow(`
		SELECT c.limit_count, 
		       (SELECT COUNT(*) FROM student_courses WHERE course_id = ?) as current_count
		FROM courses c
		WHERE c.id = ?
	`, req.CourseID, req.CourseID).Scan(&limitCount, &currentCount)

	if err != nil {
		log.Printf("检查课程人数失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	if currentCount >= limitCount {
		response := StudentJoinResponse{
			Code:    40002,
			Message: "课程人数已满",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 插入选课记录
	_, err = db.Exec(
		"INSERT INTO student_courses (stuId, course_id) VALUES (?, ?)",
		req.StuID, req.CourseID,
	)

	if err != nil {
		log.Printf("插入选课记录失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 成功响应
	response := StudentJoinResponse{
		Code:    20000,
		Message: "加入课程成功",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// myCourses 获取我的课程的核心逻辑
func myCourses(w http.ResponseWriter, stuID string, page, size int) {
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	offset := (page - 1) * size

	// 查询我的课程
	query := `
		SELECT c.*, t.name as teacher_name, t.career as teacher_career
		FROM student_courses sc
		JOIN courses c ON sc.course_id = c.id
		LEFT JOIN teacher_courses tc ON c.id = tc.course_id
		LEFT JOIN teachers t ON tc.teacher_id = t.id
		WHERE sc.stuId = ?
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, stuID, size, offset)
	if err != nil {
		log.Printf("查询我的课程失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}
	defer rows.Close()

	var courses []MyCourse
	for rows.Next() {
		var course struct {
			ID           int
			Title        string
			Cover        string
			LessonNum    int
			LimitCount   int
			TeacherName  sql.NullString
			TeacherCareer sql.NullString
		}

		// 跳过不需要的字段：description, credit, created_at, updated_at
		err := rows.Scan(
			&course.ID,
			&course.Title,
			new(interface{}), // description - 跳过
			&course.Cover,
			&course.LessonNum,
			new(interface{}), // credit - 跳过
			&course.LimitCount,
			new(interface{}), // created_at - 跳过
			new(interface{}), // updated_at - 跳过
			&course.TeacherName,
			&course.TeacherCareer,
		)

		if err != nil {
			log.Printf("解析课程数据失败: %v", err)
			continue
		}

		myCourse := MyCourse{
			ID:          course.ID,
			ClassName:   course.Title,
			TeacherName: course.TeacherName.String,
			Count1:      course.LessonNum,
			Count2:      course.LimitCount,
			Cover:       course.Cover,
			Career:      course.TeacherCareer.String,
		}

		if myCourse.TeacherName == "" {
			myCourse.TeacherName = "未知教师"
		}

		courses = append(courses, myCourse)
	}

	// 获取总数
	var total int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM student_courses WHERE stuId = ?",
		stuID,
	).Scan(&total)

	if err != nil {
		log.Printf("查询课程总数失败: %v", err)
		total = len(courses) // 如果查询失败，使用当前页的数量
	}

	// 构建响应
	response := MyCoursesResponse{
		Code: 20000,
	}
	response.Data.Courses = courses
	response.Data.Total = total

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// checkEnrollment 检查选课状态的核心逻辑
func checkEnrollment(w http.ResponseWriter, courseID int, stuID string) {
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 检查是否已选课
	var isEnrolled bool
	err = db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM student_courses WHERE stuId = ? AND course_id = ?)",
		stuID, courseID,
	).Scan(&isEnrolled)

	if err != nil {
		log.Printf("检查选课状态失败: %v", err)
		sendStudentError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 构建响应
	response := CheckEnrollmentResponse{
		Code: 20000,
	}
	response.Data.IsEnrolled = isEnrolled

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}