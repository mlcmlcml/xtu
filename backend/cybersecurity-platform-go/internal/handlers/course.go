package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"cybersecurity-platform-go/internal/database"
)

// Course 课程结构体
type Course struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Cover       string  `json:"cover"`
	LessonNum   int     `json:"lessonNum"`
	Credit      float64 `json:"credit"`
	LimitCount  int     `json:"limitCount"`
}

// CourseListResponse 课程列表响应
type CourseListResponse struct {
	Code int                 `json:"code"`
	Data CourseListData      `json:"data"`
}

type CourseListData struct {
	CourseList  []Course          `json:"courseList"`
	Total       int               `json:"total"`
	HotList     []string          `json:"hotList"`
	SubjectList []interface{}     `json:"subjectList"`
}

// CourseDetailResponse 课程详情响应
type CourseDetailResponse struct {
	Code int                   `json:"code"`
	Data CourseDetailData      `json:"data"`
}

type CourseDetailData struct {
	Course CourseDetail `json:"course"`
}

// CourseDetail 课程详情
type CourseDetail struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Cover       string     `json:"cover"`
	LessonNum   int        `json:"lessonNum"`
	Credit      float64    `json:"credit"`
	LimitCount  int        `json:"limitCount"`
	Chapter     []Chapter  `json:"chapter"`
	Teacher     TeacherInfo `json:"teacher"`
}

// Chapter 章节结构
type Chapter struct {
	ID       int        `json:"id"`
	Title    string     `json:"title"`
	State    int        `json:"state"`
	Children []Lesson   `json:"children"`
}

// Lesson 课时结构
type Lesson struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	VideoSourceID string `json:"videoSourceId"`
}

// TeacherInfo 教师信息
type TeacherInfo struct {
	TeacherID   int    `json:"teacherId"`
	TeacherName string `json:"teacherName"`
	Career      string `json:"career"`
	Intro       string `json:"intro"`
}

// GetCourseList 获取课程列表
func GetCourseList(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// 只处理GET请求
	if r.Method != http.MethodGet {
		sendCourseError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		return
	}
	
	// 获取查询参数
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	if pageSize < 1 {
		pageSize = 12
	}
	
	title := strings.TrimSpace(query.Get("title"))
	order, _ := strconv.Atoi(query.Get("order"))
	
	// 计算偏移量
	offset := (page - 1) * pageSize
	
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendCourseError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		return
	}
	
	// 构建查询
	var courses []Course
	var total int
	
	// 基础查询
	baseQuery := `
		SELECT 
			c.id,
			c.title,
			c.description,
			c.cover,
			c.lesson_num,
			c.credit,
			c.limit_count
		FROM courses c
	`
	
	// 添加搜索条件
	var whereClauses []string
	var params []interface{}
	
	if title != "" {
		whereClauses = append(whereClauses, "c.title LIKE ?")
		params = append(params, "%"+title+"%")
	}
	
	// 构建WHERE子句
	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}
	
	// 构建ORDER BY子句
	orderSQL := " ORDER BY c.id DESC" // 默认排序
	switch order {
	case 1:
		orderSQL = " ORDER BY c.id DESC" // 最新
	case 2:
		orderSQL = " ORDER BY c.lesson_num DESC" // 最热
	}
	
	// 查询课程列表
	querySQL := baseQuery + whereSQL + " GROUP BY c.id" + orderSQL + " LIMIT ? OFFSET ?"
	params = append(params, pageSize, offset)
	
	rows, err := db.Query(querySQL, params...)
	if err != nil {
		log.Printf("查询课程列表失败: %v", err)
		sendCourseError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		return
	}
	defer rows.Close()
	
	// 解析结果
	for rows.Next() {
		var course Course
		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.Cover,
			&course.LessonNum,
			&course.Credit,
			&course.LimitCount,
		)
		if err != nil {
			log.Printf("解析课程数据失败: %v", err)
			continue
		}
		courses = append(courses, course)
	}
	
	// 查询总数
	countSQL := "SELECT COUNT(DISTINCT c.id) FROM courses c" + whereSQL
	var countParams []interface{}
	if len(params) > 2 { // 移除LIMIT和OFFSET参数
		countParams = params[:len(params)-2]
	}
	
	err = db.QueryRow(countSQL, countParams...).Scan(&total)
	if err != nil {
		log.Printf("查询课程总数失败: %v", err)
		total = len(courses) // 如果查询失败，使用当前页的数量
	}
	
	// 热门搜索列表
	hotList := []string{"网络安全", "渗透测试", "数据加密"}
	
	// 构建响应
	response := CourseListResponse{
		Code: 20000,
		Data: CourseListData{
			CourseList:  courses,
			Total:       total,
			HotList:     hotList,
			SubjectList: []interface{}{}, // 空数组，与原API保持一致
		},
	}
	
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
	}
}

// GetCourseDetail 获取课程详情
func GetCourseDetail(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// 只处理GET请求
	if r.Method != http.MethodGet {
		sendCourseError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		return
	}
	
	// 从URL路径获取课程ID
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	
	if len(parts) < 3 {
		sendCourseError(w, http.StatusBadRequest, 400, "无效的URL路径")
		return
	}
	
	courseIDStr := parts[2]
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		sendCourseError(w, http.StatusBadRequest, 400, "无效的课程ID")
		return
	}
	
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendCourseError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		return
	}
	
	// 查询课程基本信息
	var courseDetail CourseDetail
	query := `
		SELECT 
			c.id,
			c.title,
			c.description,
			c.cover,
			c.lesson_num,
			c.credit,
			c.limit_count,
			t.id,
			t.name,
			t.career,
			t.intro
		FROM courses c
		LEFT JOIN teacher_courses tc ON c.id = tc.course_id
		LEFT JOIN teachers t ON tc.teacher_id = t.id
		WHERE c.id = ?
	`
	
	err = db.QueryRow(query, courseID).Scan(
		&courseDetail.ID,
		&courseDetail.Title,
		&courseDetail.Description,
		&courseDetail.Cover,
		&courseDetail.LessonNum,
		&courseDetail.Credit,
		&courseDetail.LimitCount,
		&courseDetail.Teacher.TeacherID,
		&courseDetail.Teacher.TeacherName,
		&courseDetail.Teacher.Career,
		&courseDetail.Teacher.Intro,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			sendCourseError(w, http.StatusNotFound, 404, "课程不存在")
		} else {
			log.Printf("查询课程详情失败: %v", err)
			sendCourseError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		}
		return
	}
	
	// 查询章节数据
	chapters, err := getCourseChapters(db, courseID)
	if err != nil {
		log.Printf("查询章节数据失败: %v", err)
		sendCourseError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		return
	}
	
	courseDetail.Chapter = chapters
	
	// 构建响应
	response := CourseDetailResponse{
		Code: 20000,
		Data: CourseDetailData{
			Course: courseDetail,
		},
	}
	
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
	}
}

// getCourseChapters 获取课程的章节数据
func getCourseChapters(db *sql.DB, courseID int) ([]Chapter, error) {
	query := `
		SELECT 
			ch.id,
			ch.title,
			ch.state,
			cc.id,
			cc.title,
			v.url
		FROM chapters ch
		LEFT JOIN chapter_children cc ON ch.id = cc.chapter_id
		LEFT JOIN videos v ON cc.video_id = v.id
		WHERE ch.course_id = ?
		ORDER BY ch.id, cc.id
	`
	
	rows, err := db.Query(query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	// 用于分组章节和课时
	chaptersMap := make(map[int]*Chapter)
	var chapters []Chapter
	
	for rows.Next() {
		var chapterID, chapterState int
		var chapterTitle string
		var childID sql.NullInt64
		var childTitle, videoURL sql.NullString
		
		err := rows.Scan(
			&chapterID,
			&chapterTitle,
			&chapterState,
			&childID,
			&childTitle,
			&videoURL,
		)
		if err != nil {
			return nil, err
		}
		
		// 检查章节是否已存在
		chapter, exists := chaptersMap[chapterID]
		if !exists {
			chapter = &Chapter{
				ID:       chapterID,
				Title:    chapterTitle,
				State:    chapterState,
				Children: []Lesson{},
			}
			chaptersMap[chapterID] = chapter
			chapters = append(chapters, *chapter)
		}
		
		// 如果有子课时，添加到章节中
		if childID.Valid && childTitle.Valid {
			lesson := Lesson{
				ID:            int(childID.Int64),
				Title:         childTitle.String,
				VideoSourceID: videoURL.String,
			}
			
			// 找到对应的章节并添加课时
			for i := range chapters {
				if chapters[i].ID == chapterID {
					chapters[i].Children = append(chapters[i].Children, lesson)
					break
				}
			}
		}
	}
	
	return chapters, nil
}

// sendCourseError 发送课程相关错误响应
func sendCourseError(w http.ResponseWriter, httpStatus, code int, message string) {
	w.WriteHeader(httpStatus)
	response := map[string]interface{}{
		"code":    code,
		"message": message,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码错误响应失败: %v", err)
	}
}

// RegisterCourseRoutes 注册课程相关路由
func RegisterCourseRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	
	// 课程列表
	mux.HandleFunc("/api/courses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetCourseList(w, r)
		} else {
			sendCourseError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		}
	})
	
	// 课程详情
	mux.HandleFunc("/api/courses/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetCourseDetail(w, r)
		} else {
			sendCourseError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		}
	})
	
	return mux
}