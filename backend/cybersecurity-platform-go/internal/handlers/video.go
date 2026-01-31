package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"cybersecurity-platform-go/internal/database"
)

// Video 视频结构体，对应数据库中的videos表
type Video struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
}

// VideoResponse 视频API响应格式
type VideoResponse struct {
	Code    int         `json:"code"`
	Data    VideoData   `json:"data"`
	Message string      `json:"message,omitempty"`
}

type VideoData struct {
	Video Video `json:"video"`
}

// ErrorResponse 错误响应格式
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetVideoByID 获取视频详情
func GetVideoByID(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// 从URL路径获取视频ID
	// 路径格式：/api/videos/123
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		sendError(w, http.StatusBadRequest, 400, "无效的URL路径")
		return
	}
	
	videoIDStr := pathParts[2] // pathParts[0]="api", [1]="videos", [2]="123"
	
	// 转换视频ID为整数
	videoID, err := strconv.Atoi(videoIDStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, 400, "无效的视频ID")
		return
	}
	
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		return
	}
	
	// 查询视频信息
	var video Video
	query := `
		SELECT 
			id,
			url,
			description,
			duration,
			created_at
		FROM videos
		WHERE id = ?
	`
	
	err = db.QueryRow(query, videoID).Scan(
		&video.ID,
		&video.URL,
		&video.Description,
		&video.Duration,
		&video.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, 404, "视频不存在")
		} else {
			log.Printf("获取视频详情失败: %v", err)
			sendError(w, http.StatusInternalServerError, 500, "服务器内部错误")
		}
		return
	}
	
	// 构建响应
	response := VideoResponse{
		Code: 20000,
		Data: VideoData{
			Video: video,
		},
	}
	
	// 记录日志（这里使用了fmt包）
	fmt.Printf("视频API调用成功：获取视频ID=%d\n", videoID)
	
	// 发送响应
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
	}
}

// sendError 发送错误响应
func sendError(w http.ResponseWriter, httpStatus, code int, message string) {
	w.WriteHeader(httpStatus)
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}
	
	// 记录错误日志（这里使用了fmt包）
	fmt.Printf("视频API错误：code=%d, message=%s\n", code, message)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码错误响应失败: %v", err)
	}
}

// RegisterVideoRoutes 注册视频相关路由
func RegisterVideoRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/videos/", func(w http.ResponseWriter, r *http.Request) {
		// 记录请求信息（这里使用了fmt包）
		fmt.Printf("收到视频API请求：%s %s\n", r.Method, r.URL.Path)
		
		// 只处理GET请求
		if r.Method != http.MethodGet {
			sendError(w, http.StatusMethodNotAllowed, 405, "方法不允许")
			return
		}
		
		// 调用获取视频详情函数
		GetVideoByID(w, r)
	})
	
	return mux
}