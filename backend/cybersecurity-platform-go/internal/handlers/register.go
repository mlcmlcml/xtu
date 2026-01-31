// internal/handlers/register.go
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"cybersecurity-platform-go/internal/config"
	"cybersecurity-platform-go/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// 注册请求结构
type RegisterRequest struct {
	StuID    string `json:"stuId"`
	Email    string `json:"email"`
	Password string `json:"password"`
	NickName string `json:"nickName"`
	Avatar   []byte `json:"avatar,omitempty"`
}

// 注册响应结构
type RegisterResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	User    interface{} `json:"user,omitempty"`
}

// RegisterRoutes 注册路由
func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", registerCorsMiddleware(registerHandler))

	return mux
}

// registerHandler 注册处理器
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// 解析 multipart/form-data
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		// 尝试解析 JSON
		handleJSONRegister(w, r)
		return
	}

	// 处理表单注册
	handleFormRegister(w, r)
}

// handleJSONRegister 处理 JSON 注册
func handleJSONRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendRegisterError(w, http.StatusBadRequest, 40000, "请求体格式错误")
		return
	}

	// 验证输入
	if req.StuID == "" || req.Email == "" || req.Password == "" || req.NickName == "" {
		sendRegisterError(w, http.StatusBadRequest, 40000, "所有字段都是必需的")
		return
	}

	// 调用注册逻辑
	registerUser(w, req, "")
}

// handleFormRegister 处理表单注册
func handleFormRegister(w http.ResponseWriter, r *http.Request) {
	stuID := r.FormValue("stuId")
	email := r.FormValue("email")
	password := r.FormValue("password")
	nickName := r.FormValue("nickName")

	// 验证输入
	if stuID == "" || email == "" || password == "" || nickName == "" {
		sendRegisterError(w, http.StatusBadRequest, 40000, "所有字段都是必需的")
		return
	}

	// 处理头像文件
	var avatarFilename string
	file, handler, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()

		// 验证文件类型
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/jpg":  true,
			"image/png":  true,
			"image/gif":  true,
		}

		if !allowedTypes[handler.Header.Get("Content-Type")] {
			sendRegisterError(w, http.StatusBadRequest, 40000, "不支持的文件类型，仅支持 JPEG、PNG、GIF")
			return
		}

		// 验证文件大小（5MB限制）
		if handler.Size > 5*1024*1024 {
			sendRegisterError(w, http.StatusBadRequest, 40000, "图片大小不能超过5MB")
			return
		}

		// 加载配置获取用户头像目录
		cfg := config.LoadConfig()
		userImageDir := cfg.UserImageDir
		if userImageDir == "" {
			userImageDir = filepath.Join("assets", "image", "user")
		}

		// 确保目录存在
		if err := os.MkdirAll(userImageDir, 0755); err != nil {
			log.Printf("创建目录失败: %v", err)
			sendRegisterError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
			return
		}

		// 生成唯一文件名
		ext := filepath.Ext(handler.Filename)
		if ext == "" {
			ext = ".jpg"
		}
		avatarFilename = fmt.Sprintf("avatar-%d%s", time.Now().UnixNano(), ext)
		avatarPath := filepath.Join(userImageDir, avatarFilename)

		// 保存文件
		dst, err := os.Create(avatarPath)
		if err != nil {
			log.Printf("创建文件失败: %v", err)
			sendRegisterError(w, http.StatusInternalServerError, 50000, "保存头像失败")
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("保存文件内容失败: %v", err)
			sendRegisterError(w, http.StatusInternalServerError, 50000, "保存头像失败")
			return
		}
	}

	// 构建注册请求
	req := RegisterRequest{
		StuID:    stuID,
		Email:    email,
		Password: password,
		NickName: nickName,
	}

	// 调用注册逻辑
	registerUser(w, req, avatarFilename)
}

// registerUser 用户注册逻辑
func registerUser(w http.ResponseWriter, req RegisterRequest, avatarFilename string) {
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendRegisterError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 检查学号是否已存在
	var stuIDExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM students WHERE stuId = ?)", req.StuID).Scan(&stuIDExists)
	if err != nil {
		log.Printf("检查学号失败: %v", err)
		sendRegisterError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	if stuIDExists {
		sendRegisterError(w, http.StatusBadRequest, 40001, "学号已被注册")
		return
	}

	// 检查邮箱是否已存在
	var emailExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM students WHERE email = ?)", req.Email).Scan(&emailExists)
	if err != nil {
		log.Printf("检查邮箱失败: %v", err)
		sendRegisterError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	if emailExists {
		sendRegisterError(w, http.StatusBadRequest, 40002, "邮箱已被注册")
		return
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		sendRegisterError(w, http.StatusInternalServerError, 50000, "密码加密失败")
		return
	}

	// 构建用户头像URL
	userHead := "https://tse3-mm.cn.bing.net/th/id/OIP-C.qidgOqAsPEdzAg5inmSK3AAAAA?rs=1&pid=ImgDetMain" // 默认头像
	if avatarFilename != "" {
		userHead = fmt.Sprintf("/img/user/%s", avatarFilename)
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		log.Printf("开启事务失败: %v", err)
		sendRegisterError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}
	defer tx.Rollback()

	// 插入到 students 表
	insertStudentSQL := "INSERT INTO students (stuId, password, email) VALUES (?, ?, ?)"
	result, err := tx.Exec(insertStudentSQL, req.StuID, string(hashedPassword), req.Email)
	if err != nil {
		log.Printf("插入学生表失败: %v", err)
		sendRegisterError(w, http.StatusInternalServerError, 50000, "注册失败")
		return
	}

	studentID, _ := result.LastInsertId()

	// 插入到 userdetail 表
	insertUserDetailSQL := "INSERT INTO userdetail (stuId, nickName, userHead, userEmail) VALUES (?, ?, ?, ?)"
	_, err = tx.Exec(insertUserDetailSQL, req.StuID, req.NickName, userHead, req.Email)
	if err != nil {
		log.Printf("插入用户详情表失败: %v", err)
		sendRegisterError(w, http.StatusInternalServerError, 50000, "注册失败")
		return
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		log.Printf("提交事务失败: %v", err)
		sendRegisterError(w, http.StatusInternalServerError, 50000, "注册失败")
		return
	}

	// 构建响应
	userInfo := map[string]interface{}{
		"id":        studentID,
		"stuId":     req.StuID,
		"userEmail": req.Email,
		"nickName":  req.NickName,
		"userHead":  userHead,
	}

	response := RegisterResponse{
		Code:    20000,
		Message: "注册成功",
		User:    userInfo,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sendRegisterError 发送注册错误响应
func sendRegisterError(w http.ResponseWriter, httpStatus, code int, message string) {
	w.WriteHeader(httpStatus)
	response := map[string]interface{}{
		"code":    code,
		"message": message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码错误响应失败: %v", err)
	}
}

// registerCorsMiddleware CORS中间件（避免与forum.go冲突）
func registerCorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}