package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cybersecurity-platform-go/internal/config"
	"cybersecurity-platform-go/internal/database"
	"cybersecurity-platform-go/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("=== 网络安全平台后端（Go版本） ===")
	fmt.Println("正在启动...")

	// 1. 加载环境变量
	envFile := ".env"
	if os.Getenv("NODE_ENV") == "production" {
		envFile = ".env.production"
		fmt.Println("检测到生产环境")
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("注意：未找到 %s 文件，使用系统环境变量", envFile)
	} else {
		fmt.Printf("✓ 已加载环境变量文件: %s\n", envFile)
	}

	// 2. 加载配置
	cfg := config.LoadConfig()

	// 显示配置信息
	fmt.Println("\n=== 配置信息 ===")
	fmt.Printf("运行环境: %s\n", cfg.Env)
	fmt.Printf("服务地址: %s\n", cfg.BaseURL)
	fmt.Printf("监听端口: %s\n", cfg.Port)
	fmt.Printf("数据库: %s@%s:%s/%s\n", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// 显示目录信息
	fmt.Println("\n=== 目录配置 ===")
	fmt.Printf("视频目录: %s\n", cfg.GetFirstExistingVideoDir())
	fmt.Printf("PDF目录: %s\n", cfg.GetFirstExistingPdfDir())
	fmt.Printf("图片目录: %s\n", cfg.GetFirstExistingImageDir())
	fmt.Printf("文章目录: %s\n", cfg.GetFirstExistingArticleDir())
	fmt.Printf("用户头像: %s\n", cfg.UserImageDir)
	fmt.Printf("课程图片: %s\n", cfg.CourseImageDir)

	// 3. 测试数据库连接
	fmt.Println("\n=== 数据库连接测试 ===")
	if err := database.TestConnection(); err != nil {
		log.Printf("⚠️  数据库连接测试失败: %v", err)
		fmt.Println("将继续启动服务，但数据库相关功能可能不可用")
	} else {
		fmt.Println("✓ 数据库连接成功")
	}

	// 4. 初始化图数据库连接
	fmt.Println("\n=== Neo4j图数据库连接测试 ===")
	if _, err := handlers.InitGraphHandler(cfg.Neo4jURI, cfg.Neo4jUser, cfg.Neo4jPassword); err != nil {
		log.Printf("⚠️  Neo4j连接失败: %v", err)
		fmt.Println("图数据库功能将不可用")
	} else {
		fmt.Println("✓ Neo4j连接成功")
		defer func() {
	if cfg.Neo4jURI != "" && cfg.Neo4jUser != "" && cfg.Neo4jPassword != "" {
		if gh := handlers.GetGraphHandler(); gh != nil {
			gh.Close()
			fmt.Println("Neo4j连接已关闭")
		}
	}
}()
	}

	// 5. 注册路由
	fmt.Println("\n=== 注册API路由 ===")

	// 5.1 登录路由
	loginMux := handlers.RegisterLoginRoutes()
	fmt.Println("✓ 登录API路由已注册: /api/login")

	// 5.2 注册路由（假设您已有）
	// registerMux := handlers.RegisterRoutes()
	// fmt.Println("✓ 注册API路由已注册: /api/register")

	// 5.3 视频路由
	videoMux := handlers.RegisterVideoRoutes()
	fmt.Println("✓ 视频API路由已注册: /api/videos")

	// 5.4 课程路由
	courseMux := handlers.RegisterCourseRoutes()
	fmt.Println("✓ 课程API路由已注册: /api/courses")

	// 5.5 学生路由
	studentMux := handlers.RegisterStudentRoutes()
	fmt.Println("✓ 学生API路由已注册: /api/student")

	// 5.6 教师路由
	teacherMux := handlers.RegisterTeacherRoutes()
	fmt.Println("✓ 教师API路由已注册: /api/teachers")

	// 5.7 论坛路由
	forumMux := handlers.RegisterForumRoutes()
	fmt.Println("✓ 论坛API路由已注册: /api/forum")

	// 5.8 图数据库路由
	graphMux := handlers.RegisterGraphRoutes()
	fmt.Println("✓ 图数据库API路由已注册: /api/init-graph, /api/expand-node")

	// 6. 创建主路由处理器
	mainMux := http.NewServeMux()

	// 7. 添加各个路由到主路由
	mainMux.Handle("/api/videos/", videoMux)
	mainMux.Handle("/api/login", loginMux)
	// mainMux.Handle("/api/register", registerMux)
	mainMux.Handle("/api/courses/", courseMux)
	mainMux.Handle("/api/student/", studentMux)
	mainMux.Handle("/api/teachers/", teacherMux)
	mainMux.Handle("/api/forum/", forumMux)
	mainMux.Handle("/api/", graphMux)

	fmt.Println("✓ 所有路由已添加到主路由")

	// 8. 添加静态文件服务（支持多个目录）
	fmt.Println("\n=== 注册静态文件服务 ===")

	// 8.1 用户头像静态服务（多个可能位置）
	registerMultiDirStatic(mainMux, "/img/user/", []string{
		cfg.UserImageDir,
		filepath.Join(cwd(), "static", "images", "user"),
		filepath.Join(cwd(), "MAIN", "server", "static", "images", "user"),
		filepath.Join(cwd(), "assets", "image", "user"),
	})
	fmt.Println("✓ 用户头像静态服务: /img/user/")

	// 8.2 课程图片静态服务
	registerMultiDirStatic(mainMux, "/img/course/", []string{
		cfg.CourseImageDir,
		filepath.Join(cwd(), "static", "images", "course"),
		filepath.Join(cwd(), "MAIN", "server", "static", "images", "course"),
		filepath.Join(cwd(), "assets", "image", "course"),
	})
	fmt.Println("✓ 课程图片静态服务: /img/course/")

	// 8.3 视频静态服务（对应原 /api/videoing）
	registerMultiDirStatic(mainMux, "/api/videoing/", cfg.VideoDirs)
	fmt.Println("✓ 视频静态服务: /api/videoing/")

	// 8.4 PDF静态服务（对应原 /api/pdfs）
	registerMultiDirStatic(mainMux, "/api/pdfs/", cfg.PdfDirs)
	fmt.Println("✓ PDF静态服务: /api/pdfs/")

	// 8.5 论坛文章内容静态服务
	registerMultiDirStatic(mainMux, "/api/forum/articles/content/", cfg.ArticleDirs)
	fmt.Println("✓ 论坛文章内容静态服务: /api/forum/articles/content/")

	// 8.6 论坛上传文件静态服务
	mainMux.Handle("/api/forum/uploads/",
		http.StripPrefix("/api/forum/uploads/", http.FileServer(http.Dir(cfg.ForumUploadDir))))
	fmt.Println("✓ 论坛上传文件静态服务: /api/forum/uploads/")

	// 8.7 通用图片静态服务（备用）
	registerMultiDirStatic(mainMux, "/images/", cfg.ImageDirs)
	fmt.Println("✓ 通用图片静态服务: /images/")

	// 9. 健康检查端点
	mainMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		// 检查各个目录状态
		videoDirStatus := checkDirStatus(cfg.GetFirstExistingVideoDir())
		pdfDirStatus := checkDirStatus(cfg.GetFirstExistingPdfDir())
		articleDirStatus := checkDirStatus(cfg.GetFirstExistingArticleDir())
		userImageDirStatus := checkDirStatus(cfg.UserImageDir)
		
		fmt.Fprintf(w, `{
			"status": "ok", 
			"service": "cybersecurity-platform-go", 
			"version": "1.0.0",
			"environment": "%s",
			"database": {
				"mysql": "%s",
				"neo4j": "%s"
			},
			"directories": {
				"videos": {"path": "%s", "exists": %t},
				"pdfs": {"path": "%s", "exists": %t},
				"articles": {"path": "%s", "exists": %t},
				"user_images": {"path": "%s", "exists": %t}
			}
		}`, 
		cfg.Env,
		databaseStatus(),
		neo4jStatus(cfg),
		cfg.GetFirstExistingVideoDir(), videoDirStatus,
		cfg.GetFirstExistingPdfDir(), pdfDirStatus,
		cfg.GetFirstExistingArticleDir(), articleDirStatus,
		cfg.UserImageDir, userImageDirStatus)
	})
	fmt.Println("✓ 健康检查路由: /health")

	// 10. 首页路由
	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		
		// 获取各个目录状态
		videoDirs := generateDirStatusList(cfg.VideoDirs)
		pdfDirs := generateDirStatusList(cfg.PdfDirs)
		articleDirs := generateDirStatusList(cfg.ArticleDirs)
		imageDirs := generateDirStatusList(cfg.ImageDirs)
		
		// Neo4j状态
		neo4jStatusHTML := `<span style="color: #dc3545;">❌ 未连接</span>`
		if cfg.Neo4jURI != "" && cfg.Neo4jUser != "" && cfg.Neo4jPassword != "" {
			if gh := handlers.GetGraphHandler(); gh != nil {
				neo4jStatusHTML = `<span style="color: #28a745;">✅ 已连接</span>`
			}
		}
		
		html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>网络安全平台后端 (Go版本)</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        
        header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }
        
        h1 {
            font-size: 2.5rem;
            margin-bottom: 10px;
            font-weight: 700;
        }
        
        .subtitle {
            font-size: 1.2rem;
            opacity: 0.9;
            margin-bottom: 20px;
        }
        
        .status-badge {
            display: inline-block;
            background: rgba(255,255,255,0.2);
            padding: 8px 20px;
            border-radius: 50px;
            font-size: 0.9rem;
            font-weight: 600;
            margin-top: 10px;
        }
        
        main {
            padding: 40px;
        }
        
        .section {
            margin-bottom: 40px;
            padding-bottom: 30px;
            border-bottom: 1px solid #eee;
        }
        
        .section:last-child {
            border-bottom: none;
            margin-bottom: 0;
            padding-bottom: 0;
        }
        
        h2 {
            font-size: 1.8rem;
            color: #667eea;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #f0f0f0;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        h2::before {
            content: '';
            display: inline-block;
            width: 6px;
            height: 24px;
            background: #667eea;
            border-radius: 3px;
        }
        
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }
        
        .card {
            background: #f8f9fa;
            border-radius: 12px;
            padding: 20px;
            transition: all 0.3s ease;
            border: 1px solid #e9ecef;
        }
        
        .card:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            border-color: #667eea;
        }
        
        .api-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 12px 0;
            border-bottom: 1px solid #e9ecef;
        }
        
        .api-item:last-child {
            border-bottom: none;
        }
        
        .method {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.8rem;
            font-weight: 600;
            margin-right: 10px;
        }
        
        .method.get {
            background: #28a745;
            color: white;
        }
        
        .method.post {
            background: #007bff;
            color: white;
        }
        
        .method.put {
            background: #ffc107;
            color: #212529;
        }
        
        .method.delete {
            background: #dc3545;
            color: white;
        }
        
        .status-indicator {
            font-size: 0.9rem;
            font-weight: 600;
            padding: 4px 12px;
            border-radius: 20px;
        }
        
        .status-active {
            background: #d4edda;
            color: #155724;
        }
        
        .status-warning {
            background: #fff3cd;
            color: #856404;
        }
        
        .dir-list {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 15px;
            margin: 10px 0;
            font-family: 'Courier New', monospace;
            font-size: 0.9rem;
        }
        
        .dir-item {
            padding: 8px;
            margin: 4px 0;
            border-radius: 4px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .dir-item:nth-child(odd) {
            background: #e9ecef;
        }
        
        .checkmark {
            color: #28a745;
            font-weight: bold;
        }
        
        .crossmark {
            color: #dc3545;
            font-weight: bold;
        }
        
        .info-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }
        
        .info-item {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        
        .info-label {
            font-size: 0.9rem;
            color: #6c757d;
            margin-bottom: 5px;
        }
        
        .info-value {
            font-size: 1.1rem;
            font-weight: 600;
            color: #495057;
        }
        
        footer {
            background: #f8f9fa;
            padding: 30px 40px;
            text-align: center;
            border-top: 1px solid #e9ecef;
            color: #6c757d;
        }
        
        .quick-links {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-top: 20px;
            flex-wrap: wrap;
        }
        
        .link-button {
            display: inline-block;
            padding: 10px 20px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            transition: all 0.3s ease;
        }
        
        .link-button:hover {
            background: #5a67d8;
            transform: translateY(-2px);
        }
        
        @media (max-width: 768px) {
            .container {
                border-radius: 10px;
            }
            
            header {
                padding: 30px 20px;
            }
            
            main {
                padding: 20px;
            }
            
            h1 {
                font-size: 2rem;
            }
            
            .grid {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🚀 网络安全平台后端</h1>
            <div class="subtitle">Go语言版本 - 高性能后端服务</div>
            <div class="status-badge">✅ 服务正在运行 - 环境: ` + cfg.Env + `</div>
        </header>
        
        <main>
            <!-- 系统信息 -->
            <div class="section">
                <h2>📊 系统信息</h2>
                <div class="info-grid">
                    <div class="info-item">
                        <div class="info-label">运行环境</div>
                        <div class="info-value">` + cfg.Env + `</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">服务地址</div>
                        <div class="info-value">` + cfg.BaseURL + `</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">监听端口</div>
                        <div class="info-value">` + cfg.Port + `</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">MySQL数据库</div>
                        <div class="info-value">` + databaseStatus() + `</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Neo4j图数据库</div>
                        <div class="info-value">` + neo4jStatusHTML + `</div>
                    </div>
                </div>
            </div>
            
            <!-- 核心API -->
            <div class="section">
                <h2>🔌 核心API接口</h2>
                <div class="grid">
                    <div class="card">
                        <h3 style="margin-bottom: 15px; color: #495057;">🔐 认证模块</h3>
                        <div class="api-item">
                            <div>
                                <span class="method post">POST</span>
                                <code>/api/login</code>
                            </div>
                            <span class="status-indicator status-active">可用</span>
                        </div>
                    </div>
                    
                    <div class="card">
                        <h3 style="margin-bottom: 15px; color: #495057;">📚 课程模块</h3>
                        <div class="api-item">
                            <div>
                                <span class="method get">GET</span>
                                <code>/api/courses</code>
                            </div>
                            <span class="status-indicator status-active">可用</span>
                        </div>
                        <div class="api-item">
                            <div>
                                <span class="method get">GET</span>
                                <code>/api/courses/{id}</code>
                            </div>
                            <span class="status-indicator status-active">可用</span>
                        </div>
                    </div>
                    
                    <div class="card">
                        <h3 style="margin-bottom: 15px; color: #495057;">💬 论坛模块</h3>
                        <div class="api-item">
                            <div>
                                <span class="method get">GET</span>
                                <code>/api/forum/articles</code>
                            </div>
                            <span class="status-indicator status-active">可用</span>
                        </div>
                        <div class="api-item">
                            <div>
                                <span class="method get">GET</span>
                                <code>/api/forum/categories</code>
                            </div>
                            <span class="status-indicator status-active">可用</span>
                        </div>
                    </div>
                    
                    <div class="card">
                        <h3 style="margin-bottom: 15px; color: #495057;">📊 图数据库</h3>
                        <div class="api-item">
                            <div>
                                <span class="method get">GET</span>
                                <a href="/api/init-graph"><code>/api/init-graph</code></a>
                            </div>
                            <span class="status-indicator">` + neo4jStatusHTML + `</span>
                        </div>
                        <div class="api-item">
                            <div>
                                <span class="method get">GET</span>
                                <code>/api/expand-node/{name}</code>
                            </div>
                            <span class="status-indicator">` + neo4jStatusHTML + `</span>
                        </div>
                    </div>
                </div>
            </div>
            
            <!-- 目录配置 -->
            <div class="section">
                <h2>📁 目录配置状态</h2>
                
                <div class="card">
                    <h3 style="margin-bottom: 15px; color: #495057;">🎬 视频目录</h3>
                    <div class="dir-list">
                        ` + videoDirs + `
                    </div>
                </div>
                
                <div class="card">
                    <h3 style="margin-bottom: 15px; color: #495057;">📄 PDF目录</h3>
                    <div class="dir-list">
                        ` + pdfDirs + `
                    </div>
                </div>
                
                <div class="card">
                    <h3 style="margin-bottom: 15px; color: #495057;">📝 文章目录</h3>
                    <div class="dir-list">
                        ` + articleDirs + `
                    </div>
                </div>
                
                <div class="card">
                    <h3 style="margin-bottom: 15px; color: #495057;">🖼️ 图片目录</h3>
                    <div class="dir-list">
                        ` + imageDirs + `
                    </div>
                </div>
            </div>
        </main>
        
        <footer>
            <p>© 2024 网络安全平台 - Go版本 | 高性能后端服务</p>
            <div class="quick-links">
                <a href="/health" class="link-button">健康检查</a>
                <a href="/api/courses" class="link-button">课程列表</a>
                <a href="/api/forum/categories" class="link-button">论坛分类</a>
                <a href="/img/user/" class="link-button">用户头像</a>
                <a href="/api/videoing/" class="link-button">视频文件</a>
            </div>
        </footer>
    </div>
</body>
</html>`
		
		fmt.Fprintf(w, html)
	})
	fmt.Println("✓ 首页路由: /")

	// 11. 设置服务器地址
	serverAddr := ":" + cfg.Port

	// 12. 显示启动信息
	fmt.Println("\n=== 启动信息 ===")
	fmt.Printf("服务器将启动在: %s\n", cfg.BaseURL)
	fmt.Printf("首页: %s\n", cfg.BaseURL)
	fmt.Printf("健康检查: %s/health\n", cfg.BaseURL)
	fmt.Printf("用户头像: %s/img/user/\n", cfg.BaseURL)
	fmt.Printf("课程图片: %s/img/course/\n", cfg.BaseURL)
	fmt.Printf("视频文件: %s/api/videoing/\n", cfg.BaseURL)
	fmt.Printf("PDF文件: %s/api/pdfs/\n", cfg.BaseURL)
	fmt.Printf("图数据库API: %s/api/init-graph\n", cfg.BaseURL)
	fmt.Printf("论坛文章: %s/api/forum/articles\n", cfg.BaseURL)

	if cfg.IsProduction() {
		fmt.Println("\n⚠️  生产环境注意事项:")
		fmt.Println("1. 确保 .env.production 文件已正确配置")
		fmt.Println("2. 数据库连接信息已加密")
		fmt.Println("3. 静态文件服务已启用")
		fmt.Println("4. 建议启用HTTPS")
		fmt.Println("5. 配置适当的CORS策略")
	}

	fmt.Println("\n按 Ctrl+C 停止服务器")

	// 13. 程序退出时清理资源
	defer func() {
		fmt.Println("\n正在清理资源...")
		database.CloseDB()
		if gh := handlers.GetGraphHandler(); gh != nil {
			gh.Close()
			fmt.Println("Neo4j连接已关闭")
		}
		fmt.Println("资源清理完成")
	}()

	// 14. 启动HTTP服务器
	fmt.Println("\n🚀 启动HTTP服务器...")
	log.Printf("服务已启动: %s", cfg.BaseURL)
	if err := http.ListenAndServe(serverAddr, mainMux); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// 获取当前工作目录
func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

// registerMultiDirStatic 注册多目录静态服务
func registerMultiDirStatic(mux *http.ServeMux, prefix string, dirs []string) {
	mux.Handle(prefix, http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// 处理OPTIONS请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// 只允许GET请求
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// 尝试每个目录，直到找到文件
		for _, dir := range dirs {
			filePath := filepath.Join(dir, r.URL.Path)
			if _, err := os.Stat(filePath); err == nil {
				http.ServeFile(w, r, filePath)
				return
			}
		}
		
		// 所有目录都没找到文件
		http.NotFound(w, r)
	})))
}

// 检查目录状态
func checkDirStatus(dir string) bool {
	if _, err := os.Stat(dir); err == nil {
		return true
	}
	return false
}

// 生成目录状态列表HTML
func generateDirStatusList(dirs []string) string {
	html := ""
	for i, dir := range dirs {
		exists := checkDirStatus(dir)
		statusIcon := "✅"
		statusText := "存在"
		if !exists {
			statusIcon = "❌"
			statusText = "不存在"
		}
		
		priority := ""
		if i == 0 {
			priority = " (优先)"
		}
		
		html += fmt.Sprintf(`<div class="dir-item">
            <span>%s</span>
            <span style="flex-grow: 1;">%s%s</span>
            <span style="color: %s;">%s</span>
        </div>`, 
		statusIcon, dir, priority, 
		getStatusColor(exists), statusText)
	}
	return html
}

func getStatusColor(exists bool) string {
	if exists {
		return "#28a745"
	}
	return "#dc3545"
}

func databaseStatus() string {
	if err := database.TestConnection(); err == nil {
		return "✅ 已连接"
	}
	return "❌ 未连接"
}

func neo4jStatus(cfg *config.Config) string {
	if cfg.Neo4jURI == "" || cfg.Neo4jUser == "" || cfg.Neo4jPassword == "" {
		return "未配置"
	}
	
	if gh := handlers.GetGraphHandler(); gh != nil {
		return "已连接"
	}
	return "连接失败"
}