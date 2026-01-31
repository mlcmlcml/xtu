package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"cybersecurity-platform-go/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求结构体
type LoginRequest struct {
	StuID    string `json:"stuId"`
	Password string `json:"password"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	User    *UserInfo   `json:"user,omitempty"`
}

// UserInfo 用户信息结构体
type UserInfo struct {
	ID        int    `json:"id"`
	StuID     string `json:"stuId"`
	UserEmail string `json:"userEmail"`
	UserHead  string `json:"userHead"`
	UserName  string `json:"userName"`
	NickName  string `json:"nickName"`
}

// LoginHandler 处理登录请求
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// 只处理POST请求
	if r.Method != http.MethodPost {
		sendLoginError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
		return
	}
	
	// 解析请求体
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendLoginError(w, http.StatusBadRequest, 40000, "请求格式错误")
		return
	}
	
	// 验证输入
	if req.StuID == "" || req.Password == "" {
		sendLoginError(w, http.StatusBadRequest, 40000, "学号和密码都是必需的")
		return
	}
	
	// 清理输入
	req.StuID = strings.TrimSpace(req.StuID)
	
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendLoginError(w, http.StatusInternalServerError, 50000, "登录失败，请重试")
		return
	}
	
	// 查询用户信息
	query := `
		SELECT 
			b.id,
			a.stuId, 
			a.password,
			a.email,
			b.userName,
			b.userHead,
			b.nickName
		FROM students a 
		JOIN userdetail b ON a.stuId = b.stuId 
		WHERE a.stuId = ?
	`
	
	var user UserInfo
	var hashedPassword string
	
	err = db.QueryRow(query, req.StuID).Scan(
		&user.ID,
		&user.StuID,
		&hashedPassword,
		&user.UserEmail,
		&user.UserName,
		&user.UserHead,
		&user.NickName,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			sendLoginError(w, http.StatusUnauthorized, 40001, "用户不存在")
		} else {
			log.Printf("数据库查询错误: %v", err)
			sendLoginError(w, http.StatusInternalServerError, 50000, "登录失败，请重试")
		}
		return
	}
	
	// 验证密码（使用bcrypt）
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		// 密码错误
		sendLoginError(w, http.StatusUnauthorized, 40002, "密码错误")
		return
	}
	
	// 登录成功
	response := LoginResponse{
		Code:    20000,
		Message: "登录成功",
		User:    &user,
	}
	
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
	}
}

// sendLoginError 发送登录错误响应
func sendLoginError(w http.ResponseWriter, httpStatus, code int, message string) {
	w.WriteHeader(httpStatus)
	response := LoginResponse{
		Code:    code,
		Message: message,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码错误响应失败: %v", err)
	}
}

// RegisterLoginRoutes 注册登录相关路由
func RegisterLoginRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	
	// 登录路由
	mux.HandleFunc("/api/login", LoginHandler)
	
	return mux
}