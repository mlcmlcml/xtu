package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"cybersecurity-platform-go/internal/database"
)

// Teacher 教师信息结构体
type Teacher struct {
	ID          int    `json:"id"`
	TeacherName string `json:"teacherName"` // 对应JS中的teacherName
	Career      string `json:"career"`
	Intro       string `json:"intro"`
}

// TeacherCourse 教师课程信息
type TeacherCourse struct {
	CourseID    int    `json:"courseId"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
}

// TeacherDetailResponse 教师详情响应
type TeacherDetailResponse struct {
	Code int `json:"code"`
	Data struct {
		Teacher Teacher         `json:"teacher"`
		Course  []TeacherCourse `json:"course"`
	} `json:"data"`
}

// TeacherListResponse 教师列表响应
type TeacherListResponse struct {
	Code int `json:"code"`
	Data struct {
		Rows  []Teacher `json:"rows"`
		Total int       `json:"total"`
	} `json:"data"`
}

// RegisterTeacherRoutes 注册教师路由
func RegisterTeacherRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// 教师详情
	mux.HandleFunc("GET /api/teachers/{id}", teacherDetailHandler)

	// 教师列表
	mux.HandleFunc("GET /api/teachers", teacherListHandler)

	return mux
}

// teacherDetailHandler 教师详情处理器
func teacherDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 只处理GET请求
	if r.Method != http.MethodGet {
		sendTeacherError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		return
	}

	// 从URL路径获取教师ID
	idStr := r.PathValue("id")
	if idStr == "" {
		sendTeacherError(w, http.StatusBadRequest, 400, "缺少教师ID")
		return
	}

	teacherID, err := strconv.Atoi(idStr)
	if err != nil || teacherID <= 0 {
		sendTeacherError(w, http.StatusBadRequest, 400, "无效的教师ID")
		return
	}

	// 调用获取教师详情逻辑
	getTeacherDetail(w, teacherID)
}

// teacherListHandler 教师列表处理器
func teacherListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 只处理GET请求
	if r.Method != http.MethodGet {
		sendTeacherError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		return
	}

	// 解析查询参数
	query := r.URL.Query()
	pageStr := query.Get("page")
	pageSizeStr := query.Get("pageSize")
	name := query.Get("name")

	// 设置默认值
	page := 1
	pageSize := 8

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// 调用获取教师列表逻辑
	getTeacherList(w, page, pageSize, name)
}

// sendTeacherError 发送教师相关错误响应
func sendTeacherError(w http.ResponseWriter, httpStatus, code int, message string) {
	w.WriteHeader(httpStatus)
	response := map[string]interface{}{
		"code":    code,
		"message": message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码错误响应失败: %v", err)
	}
}

// getTeacherDetail 获取教师详情的核心逻辑
func getTeacherDetail(w http.ResponseWriter, teacherID int) {
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendTeacherError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		return
	}

	// 查询教师基本信息
	var teacher Teacher
	query := `
		SELECT 
			t.id,
			t.name AS teacherName,
			t.career,
			t.intro
		FROM teachers t
		WHERE t.id = ?
	`
	
	err = db.QueryRow(query, teacherID).Scan(
		&teacher.ID,
		&teacher.TeacherName,
		&teacher.Career,
		&teacher.Intro,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			sendTeacherError(w, http.StatusNotFound, 404, "教师不存在")
		} else {
			log.Printf("查询教师信息失败: %v", err)
			sendTeacherError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		}
		return
	}

	// 查询教师教授的课程
	coursesQuery := `
		SELECT 
			c.id AS courseId,
			c.title,
			c.cover,
			c.description
		FROM courses c
		INNER JOIN teacher_courses tc ON c.id = tc.course_id
		WHERE tc.teacher_id = ?
	`
	
	rows, err := db.Query(coursesQuery, teacherID)
	if err != nil {
		log.Printf("查询教师课程失败: %v", err)
		sendTeacherError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		return
	}
	defer rows.Close()

	var courses []TeacherCourse
	for rows.Next() {
		var course TeacherCourse
		err := rows.Scan(
			&course.CourseID,
			&course.Title,
			&course.Cover,
			&course.Description,
		)
		if err != nil {
			log.Printf("解析课程数据失败: %v", err)
			continue
		}
		courses = append(courses, course)
	}

	// 构建响应
	response := TeacherDetailResponse{
		Code: 20000,
	}
	response.Data.Teacher = teacher
	response.Data.Course = courses

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getTeacherList 获取教师列表的核心逻辑
func getTeacherList(w http.ResponseWriter, page, pageSize int, name string) {
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendTeacherError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		return
	}

	offset := (page - 1) * pageSize

	// 构建查询
	var whereClause string
	var params []interface{}

	if name != "" {
		whereClause = " WHERE t.name LIKE ?"
		params = append(params, "%"+name+"%")
	}

	// 查询教师列表
	query := `
		SELECT 
			t.id,
			t.name AS teacherName,
			t.career,
			t.intro
		FROM teachers t
	` + whereClause + `
		LIMIT ? OFFSET ?
	`
	
	params = append(params, pageSize, offset)
	
	rows, err := db.Query(query, params...)
	if err != nil {
		log.Printf("查询教师列表失败: %v", err)
		sendTeacherError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		return
	}
	defer rows.Close()

	var teachers []Teacher
	for rows.Next() {
		var teacher Teacher
		err := rows.Scan(
			&teacher.ID,
			&teacher.TeacherName,
			&teacher.Career,
			&teacher.Intro,
		)
		if err != nil {
			log.Printf("解析教师数据失败: %v", err)
			continue
		}
		teachers = append(teachers, teacher)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM teachers t" + whereClause
	var totalParams []interface{}
	if name != "" {
		totalParams = append(totalParams, "%"+name+"%")
	}
	
	var total int
	err = db.QueryRow(countQuery, totalParams...).Scan(&total)
	if err != nil {
		log.Printf("查询教师总数失败: %v", err)
		total = len(teachers) // 如果查询失败，使用当前页的数量
	}

	// 构建响应
	response := TeacherListResponse{
		Code: 20000,
	}
	response.Data.Rows = teachers
	response.Data.Total = total

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}