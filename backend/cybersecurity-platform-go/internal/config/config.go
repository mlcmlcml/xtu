// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config 应用配置
type Config struct {
	Env          string
	BaseURL      string
	Port         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	
	// 多个可能的静态文件目录（按照优先级）
	VideoDirs    []string // 视频目录路径数组
	PdfDirs      []string // PDF目录路径数组
	ImageDirs    []string // 图片目录路径数组
	ArticleDirs  []string // 文章目录路径数组
	
	// 特定用途的目录
	UserImageDir   string // 用户头像目录
	CourseImageDir string // 课程图片目录
	ForumUploadDir string // 论坛上传目录
	
	// Neo4j配置
	Neo4jURI     string
	Neo4jUser    string
	Neo4jPassword string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 确定环境
	env := os.Getenv("NODE_ENV")
	if env == "" {
		env = "development"
	}
	
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	
	// 构建各个目录的多个可能路径（按优先级排序）
	videoDirs := []string{
		filepath.Join(cwd, "static", "videos"),      // static/videos
		filepath.Join(cwd, "MAIN", "server", "static", "videos"), // MAIN/server/static/videos
		filepath.Join(cwd, "assets", "video"),       // assets/video
	}
	
	pdfDirs := []string{
		filepath.Join(cwd, "static", "pdfs"),        // static/pdfs
		filepath.Join(cwd, "MAIN", "server", "static", "pdfs"), // MAIN/server/static/pdfs
		filepath.Join(cwd, "assets", "pdf"),         // assets/pdf
	}
	
	imageDirs := []string{
		filepath.Join(cwd, "static", "images"),      // static/images
		filepath.Join(cwd, "MAIN", "server", "static", "images"), // MAIN/server/static/images
		filepath.Join(cwd, "assets", "image"),       // assets/image
	}
	
	articleDirs := []string{
		filepath.Join(cwd, "static", "forum", "articles"),      // static/forum/articles
		filepath.Join(cwd, "MAIN", "server", "static", "forum", "articles"), // MAIN/server/static/forum/articles
		filepath.Join(cwd, "assets", "forum", "articles"),      // assets/forum/articles
	}
	
	// 特定用途目录（使用第一个存在的路径）
	userImageDir := getFirstExistingDir([]string{
		filepath.Join(cwd, "static", "images", "user"),
		filepath.Join(cwd, "MAIN", "server", "static", "images", "user"),
		filepath.Join(cwd, "assets", "image", "user"),
	})
	
	courseImageDir := getFirstExistingDir([]string{
		filepath.Join(cwd, "static", "images", "course"),
		filepath.Join(cwd, "MAIN", "server", "static", "images", "course"),
		filepath.Join(cwd, "assets", "image", "course"),
	})
	
	forumUploadDir := filepath.Join(cwd, "static", "forum", "uploads")
	
	// 确保必要目录存在
	ensureDirExists(userImageDir)
	ensureDirExists(courseImageDir)
	ensureDirExists(forumUploadDir)
	
	return &Config{
		Env:           env,
		BaseURL:       getEnvOrDefault("BASE_URL", "http://localhost:3000"),
		Port:          getEnvOrDefault("PORT", "3000"),
		DBHost:        getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:        getEnvOrDefault("DB_PORT", "3306"),
		DBUser:        getEnvOrDefault("DB_USER", "root"),
		DBPassword:    getEnvOrDefault("DB_PASSWORD", "219332"),
		DBName:        getEnvOrDefault("DB_NAME", "cybersecurity-platform"),
		VideoDirs:     videoDirs,
		PdfDirs:       pdfDirs,
		ImageDirs:     imageDirs,
		ArticleDirs:   articleDirs,
		UserImageDir:  userImageDir,
		CourseImageDir: courseImageDir,
		ForumUploadDir: forumUploadDir,
		Neo4jURI:      getEnvOrDefault("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:     getEnvOrDefault("NEO4J_USER", "neo4j"),
		Neo4jPassword: getEnvOrDefault("NEO4J_PASSWORD", "hukaile5206"),
	}
}

// getFirstExistingDir 获取第一个存在的目录
func getFirstExistingDir(dirs []string) string {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}
	// 如果都不存在，返回第一个
	if len(dirs) > 0 {
		return dirs[0]
	}
	return "."
}

// getEnvOrDefault 获取环境变量或返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ensureDirExists 确保目录存在
func ensureDirExists(dir string) {
	if dir == "" || dir == "." {
		return
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("警告：创建目录 %s 失败: %v\n", dir, err)
		} else {
			fmt.Printf("目录 %s 创建成功\n", dir)
		}
	}
}

// IsProduction 检查是否为生产环境
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Env) == "production"
}

// GetDBDSN 获取数据库DSN连接字符串
func (c *Config) GetDBDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// GetFirstExistingVideoDir 获取第一个存在的视频目录
func (c *Config) GetFirstExistingVideoDir() string {
	return getFirstExistingDir(c.VideoDirs)
}

// GetFirstExistingPdfDir 获取第一个存在的PDF目录
func (c *Config) GetFirstExistingPdfDir() string {
	return getFirstExistingDir(c.PdfDirs)
}

// GetFirstExistingImageDir 获取第一个存在的图片目录
func (c *Config) GetFirstExistingImageDir() string {
	return getFirstExistingDir(c.ImageDirs)
}

// GetFirstExistingArticleDir 获取第一个存在的文章目录
func (c *Config) GetFirstExistingArticleDir() string {
	return getFirstExistingDir(c.ArticleDirs)
}