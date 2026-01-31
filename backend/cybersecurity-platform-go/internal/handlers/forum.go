// internal/handlers/forum.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cybersecurity-platform-go/internal/database"

	"github.com/patrickmn/go-cache"
)

// 创建内存缓存，过期时间5分钟，清理间隔10分钟
var forumCache = cache.New(5*time.Minute, 10*time.Minute)

// ForumCategory 论坛分类
type ForumCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ForumTag 论坛标签
type ForumTag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Hot  int    `json:"hot,omitempty"`
}

// ForumArticle 论坛文章
type ForumArticle struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	StuID      string     `json:"stuId"`
	StuName    string     `json:"stuName"`
	StuHead    string     `json:"stuHead"`
	CateID     int        `json:"cateId"`
	IsTop      bool       `json:"isTop"`
	IsEss      bool       `json:"isEss"`
	ViewCount  int        `json:"viewCount"`
	CreateTime string     `json:"createTime"`
	UpdateTime string     `json:"updateTime"`
	TagList    []ForumTag `json:"tagList,omitempty"`
}

// ForumArticleDetail 文章详情
type ForumArticleDetail struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	StuID      string     `json:"stuId"`
	StuName    string     `json:"stuName"`
	StuHead    string     `json:"stuHead"`
	CateID     int        `json:"cateId"`
	CateName   string     `json:"cateName,omitempty"`
	IsTop      bool       `json:"isTop"`
	IsEss      bool       `json:"isEss"`
	ViewCount  int        `json:"viewCount"`
	CreateTime string     `json:"createTime"`
	UpdateTime string     `json:"updateTime"`
	Content    string     `json:"content"`
	TagList    []ForumTag `json:"tagList"`
}

// ForumComment 论坛评论
type ForumComment struct {
	ID         int    `json:"id"`
	ArticleID  int    `json:"article_id"`
	Content    string `json:"content"`
	AuthorID   string `json:"author_id"`
	AuthorName string `json:"author_name"`
	AuthorHead string `json:"author_head"`
	ParentID   *int   `json:"parent_id"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
	LikeCount  int    `json:"like_count"`
}

// BaseResponse 基础响应
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// RegisterForumRoutes 注册论坛路由
func RegisterForumRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// 分类列表
	mux.HandleFunc("GET /api/forum/categories", corsMiddleware(getCategoriesHandler))

	// 文章列表
	mux.HandleFunc("GET /api/forum/articles", corsMiddleware(getArticlesHandler))

	// 热门文章
	mux.HandleFunc("GET /api/forum/articles/hot", corsMiddleware(getHotArticlesHandler))

	// 文章详情
	mux.HandleFunc("GET /api/forum/articles/{id}", corsMiddleware(getArticleDetailHandler))

	// 热门标签
	mux.HandleFunc("GET /api/forum/tags/hot", corsMiddleware(getHotTagsHandler))

	// 评论列表
	mux.HandleFunc("GET /api/forum/comments", corsMiddleware(getCommentsHandler))

	// 发表评论
	mux.HandleFunc("POST /api/forum/comments", corsMiddleware(postCommentHandler))

	// 删除评论
	mux.HandleFunc("DELETE /api/forum/comments/{id}", corsMiddleware(deleteCommentHandler))

	// 点赞评论
	mux.HandleFunc("POST /api/forum/comments/{id}/like", corsMiddleware(likeCommentHandler))

	return mux
}

// corsMiddleware CORS中间件
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

// getCategoriesHandler 获取分类列表处理器
func getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// 检查缓存
	cacheKey := "forum_categories"
	if cached, found := forumCache.Get(cacheKey); found {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cached)
		return
	}

	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	query := "SELECT id, name FROM forum_categories ORDER BY id"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("查询分类列表失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}
	defer rows.Close()

	var categories []ForumCategory
	categories = append(categories, ForumCategory{ID: 0, Name: "全部"})

	for rows.Next() {
		var category ForumCategory
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			log.Printf("解析分类数据失败: %v", err)
			continue
		}
		categories = append(categories, category)
	}

	response := BaseResponse{
		Code:    20000,
		Message: "成功获取分类列表",
		Data:    categories,
	}

	// 存入缓存
	forumCache.Set(cacheKey, response, cache.DefaultExpiration)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getArticlesHandler 获取文章列表处理器
func getArticlesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	cateId, _ := strconv.Atoi(query.Get("cateId"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建缓存键
	cacheKey := fmt.Sprintf("articles_%d_%d_%d", cateId, page, pageSize)
	if cached, found := forumCache.Get(cacheKey); found {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cached)
		return
	}

	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 构建查询
	var whereClauses []string
	var params []interface{}

	if cateId != 0 {
		whereClauses = append(whereClauses, "a.cateId = ?")
		params = append(params, cateId)
	}

	baseQuery := `
		SELECT 
			a.id, a.title, a.stuId, a.stuName, a.stuHead, 
			a.cateId, a.isTop, a.isEss, a.viewCount, 
			a.createTime, a.updateTime
		FROM forum_articles a
	`

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	baseQuery += " ORDER BY a.isTop DESC, a.updateTime DESC"
	baseQuery += " LIMIT ? OFFSET ?"
	params = append(params, pageSize, offset)

	rows, err := db.Query(baseQuery, params...)
	if err != nil {
		log.Printf("查询文章列表失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}
	defer rows.Close()

	var articles []ForumArticle
	for rows.Next() {
		var article ForumArticle
		err := rows.Scan(
			&article.ID, &article.Title, &article.StuID, &article.StuName, &article.StuHead,
			&article.CateID, &article.IsTop, &article.IsEss, &article.ViewCount,
			&article.CreateTime, &article.UpdateTime,
		)
		if err != nil {
			log.Printf("解析文章数据失败: %v", err)
			continue
		}

		articles = append(articles, article)
	}

	// 获取总数
	countQuery := "SELECT COUNT(*) AS total FROM forum_articles a"
	if len(whereClauses) > 0 {
		countQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var total int
	err = db.QueryRow(countQuery, params[:len(params)-2]...).Scan(&total)
	if err != nil {
		log.Printf("查询文章总数失败: %v", err)
		total = len(articles)
	}

	response := BaseResponse{
		Code: 20000,
		Data: map[string]interface{}{
			"items": articles,
			"total": total,
		},
	}

	// 存入缓存
	forumCache.Set(cacheKey, response, cache.DefaultExpiration)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getHotArticlesHandler 获取热门文章处理器
func getHotArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// 检查缓存
	cacheKey := "hot_articles"
	if cached, found := forumCache.Get(cacheKey); found {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cached)
		return
	}

	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	query := `
		SELECT id, title, viewCount, createTime
		FROM forum_articles 
		WHERE viewCount > 0
		ORDER BY viewCount DESC 
		LIMIT 5
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("查询热门文章失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}
	defer rows.Close()

	var hotArticles []map[string]interface{}
	for rows.Next() {
		var id, viewCount int
		var title, createTime string

		err := rows.Scan(&id, &title, &viewCount, &createTime)
		if err != nil {
			log.Printf("解析热门文章数据失败: %v", err)
			continue
		}

		// 格式化时间
		createTimeFormatted := createTime
		if t, err := time.Parse("2006-01-02 15:04:05", createTime); err == nil {
			createTimeFormatted = t.Format("2006-01-02 15:04")
		}

		hotArticles = append(hotArticles, map[string]interface{}{
			"id":         id,
			"title":      title,
			"viewCount":  viewCount,
			"createTime": createTimeFormatted,
		})
	}

	var response BaseResponse
	if len(hotArticles) == 0 {
		response = BaseResponse{
			Code:    20000,
			Message: "暂无热门文章",
			Data:    []interface{}{},
		}
	} else {
		response = BaseResponse{
			Code:    20000,
			Message: "热门文章获取成功",
			Data:    hotArticles,
		}
	}

	// 存入缓存
	forumCache.Set(cacheKey, response, cache.DefaultExpiration)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getArticleDetailHandler 获取文章详情处理器
func getArticleDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendForumError(w, http.StatusBadRequest, 40000, "无效的文章ID")
		return
	}

	// 检查缓存
	cacheKey := fmt.Sprintf("article_%d", id)
	if cached, found := forumCache.Get(cacheKey); found {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cached)
		return
	}

	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 异步更新浏览数
	go func() {
		dbLocal, err := database.GetDB()
		if err == nil {
			defer dbLocal.Close()
			_, _ = dbLocal.Exec("UPDATE forum_articles SET viewCount = viewCount + 1 WHERE id = ?", id)
		}
	}()

	// 获取文章基本信息及分类名称
	var article ForumArticleDetail
	query := `
		SELECT a.*, c.name AS cateName 
		FROM forum_articles a
		LEFT JOIN forum_categories c ON a.cateId = c.id
		WHERE a.id = ?
	`

	err = db.QueryRow(query, id).Scan(
		&article.ID, &article.Title, &article.StuID, &article.StuName, &article.StuHead,
		&article.CateID, &article.IsTop, &article.IsEss, &article.ViewCount,
		&article.CreateTime, &article.UpdateTime, &article.CateName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			sendForumError(w, http.StatusNotFound, 40400, "文章不存在")
		} else {
			log.Printf("查询文章详情失败: %v", err)
			sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		}
		return
	}

	// 模拟文章内容
	article.Content = fmt.Sprintf("这是文章《%s》的详细内容。\n\n作者：%s\n发布时间：%s\n\n文章ID：%d\n分类：%s\n浏览量：%d",
		article.Title, article.StuName, article.CreateTime, article.ID, article.CateName, article.ViewCount)

	response := BaseResponse{
		Code: 20000,
		Data: map[string]interface{}{
			"aclInfo": article,
		},
	}

	// 存入缓存（5分钟）
	forumCache.Set(cacheKey, response, 5*time.Minute)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getHotTagsHandler 获取热门标签处理器
func getHotTagsHandler(w http.ResponseWriter, r *http.Request) {
	// 检查缓存
	cacheKey := "hot_tags"
	if cached, found := forumCache.Get(cacheKey); found {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cached)
		return
	}

	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	query := `
		SELECT t.id, t.name, COUNT(at.article_id) as hot
		FROM forum_tags t
		JOIN article_tags at ON t.id = at.tag_id
		GROUP BY t.id
		ORDER BY hot DESC
		LIMIT 5
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("查询热门标签失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}
	defer rows.Close()

	var tags []ForumTag
	for rows.Next() {
		var tag ForumTag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Hot)
		if err != nil {
			log.Printf("解析标签数据失败: %v", err)
			continue
		}
		tags = append(tags, tag)
	}

	response := BaseResponse{
		Code:    20000,
		Message: "成功获取热门标签",
		Data:    tags,
	}

	// 存入缓存
	forumCache.Set(cacheKey, response, cache.DefaultExpiration)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getCommentsHandler 获取评论列表处理器
func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	articleId := r.URL.Query().Get("articleId")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if articleId == "" {
		sendForumError(w, http.StatusBadRequest, 40000, "文章ID不能为空")
		return
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建缓存键
	cacheKey := fmt.Sprintf("comments_%s_%d_%d", articleId, page, pageSize)
	if cached, found := forumCache.Get(cacheKey); found {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cached)
		return
	}

	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 获取评论列表
	commentQuery := `
		SELECT c.* 
		FROM forum_comments c 
		WHERE c.article_id = ? AND c.status = 1 
		ORDER BY c.create_time DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(commentQuery, articleId, pageSize, offset)
	if err != nil {
		log.Printf("查询评论列表失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}
	defer rows.Close()

	var comments []ForumComment
	for rows.Next() {
		var comment ForumComment
		err := rows.Scan(
			&comment.ID, &comment.ArticleID, &comment.Content,
			&comment.AuthorID, &comment.AuthorName, &comment.AuthorHead,
			&comment.ParentID, &comment.Status, &comment.CreateTime,
			&comment.LikeCount,
		)
		if err != nil {
			log.Printf("解析评论数据失败: %v", err)
			continue
		}
		comments = append(comments, comment)
	}

	// 获取评论总数
	var total int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM forum_comments WHERE article_id = ? AND status = 1",
		articleId,
	).Scan(&total)
	if err != nil {
		log.Printf("查询评论总数失败: %v", err)
		total = len(comments)
	}

	response := BaseResponse{
		Code: 20000,
		Data: comments,
	}

	// 存入缓存（评论缓存时间减半）
	forumCache.Set(cacheKey, response, 2*time.Minute)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// postCommentHandler 发表评论处理器
func postCommentHandler(w http.ResponseWriter, r *http.Request) {
	var comment struct {
		ArticleID  int    `json:"articleId"`
		Content    string `json:"content"`
		AuthorID   string `json:"authorId"`
		AuthorName string `json:"authorName"`
		AuthorHead string `json:"authorHead"`
		ParentID   *int   `json:"parentId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendForumError(w, http.StatusBadRequest, 40000, "请求体为空")
		return
	}

	if comment.ArticleID == 0 || comment.Content == "" || comment.AuthorID == "" || comment.AuthorName == "" {
		sendForumError(w, http.StatusBadRequest, 40000, "参数不完整")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 验证文章是否存在
	var articleExists bool

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM forum_articles WHERE id = ?)", comment.ArticleID).Scan(&articleExists)
	if err != nil {
		log.Printf("验证文章失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	if !articleExists {
		sendForumError(w, http.StatusNotFound, 40400, "文章不存在")
		return
	}

	// 插入评论
	insertQuery := `
		INSERT INTO forum_comments 
		(article_id, content, author_id, author_name, author_head, parent_id) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(insertQuery,
		comment.ArticleID, comment.Content, comment.AuthorID,
		comment.AuthorName, comment.AuthorHead, comment.ParentID,
	)
	if err != nil {
		log.Printf("插入评论失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "发表评论失败")
		return
	}

	commentID, _ := result.LastInsertId()

	// 清除相关缓存
	forumCache.Delete(fmt.Sprintf("comments_%d_1_10", comment.ArticleID))
	forumCache.Delete(fmt.Sprintf("article_%d", comment.ArticleID))

	response := BaseResponse{
		Code:    20000,
		Message: "评论发表成功",
		Data: map[string]interface{}{
			"id": commentID,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// deleteCommentHandler 删除评论处理器
func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendForumError(w, http.StatusBadRequest, 40000, "无效的评论ID")
		return
	}

	var reqBody struct {
		AuthorID  string `json:"authorId"`
		ArticleID int    `json:"articleId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		sendForumError(w, http.StatusBadRequest, 40000, "请求体为空")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 验证评论是否存在且属于该用户
	var commentExists bool
	err = db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM forum_comments WHERE id = ? AND author_id = ?)",
		id, reqBody.AuthorID,
	).Scan(&commentExists)

	if err != nil {
		log.Printf("验证评论失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	if !commentExists {
		sendForumError(w, http.StatusNotFound, 40400, "评论不存在或无权操作")
		return
	}

	// 软删除评论
	_, err = db.Exec("UPDATE forum_comments SET status = 0 WHERE id = ?", id)
	if err != nil {
		log.Printf("删除评论失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "删除评论失败")
		return
	}

	// 清除相关缓存
	if reqBody.ArticleID > 0 {
		forumCache.Delete(fmt.Sprintf("comments_%d_1_10", reqBody.ArticleID))
		forumCache.Delete(fmt.Sprintf("article_%d", reqBody.ArticleID))
	}

	response := BaseResponse{
		Code:    20000,
		Message: "评论删除成功",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// likeCommentHandler 点赞评论处理器
func likeCommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendForumError(w, http.StatusBadRequest, 40000, "无效的评论ID")
		return
	}

	var reqBody struct {
		ArticleID int `json:"articleId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		sendForumError(w, http.StatusBadRequest, 40000, "请求体为空")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}

	// 更新点赞数
	_, err = db.Exec("UPDATE forum_comments SET like_count = like_count + 1 WHERE id = ?", id)
	if err != nil {
		log.Printf("点赞评论失败: %v", err)
		sendForumError(w, http.StatusInternalServerError, 50000, "点赞失败")
		return
	}

	// 清除相关缓存
	if reqBody.ArticleID > 0 {
		forumCache.Delete(fmt.Sprintf("comments_%d_1_10", reqBody.ArticleID))
	}

	response := BaseResponse{
		Code:    20000,
		Message: "点赞成功",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sendForumError 发送论坛相关错误响应
func sendForumError(w http.ResponseWriter, httpStatus, code int, message string) {
	w.WriteHeader(httpStatus)
	response := BaseResponse{
		Code:    code,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码错误响应失败: %v", err)
	}
}