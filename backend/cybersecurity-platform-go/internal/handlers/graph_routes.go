package handlers

import (
	"encoding/json"
	"net/http"
)

// InitGraphHandler 初始化图处理器
func InitGraphHandler(uri, username, password string) (*GraphHandler, error) {
	return NewGraphHandler(uri, username, password)
}

// GetGraphHandler 获取图处理器实例
func GetGraphHandler() *GraphHandler {
	// 这里需要你实现单例模式或从某个地方获取实例
	// 暂时返回nil，你需要根据实际情况修改
	return nil
}

// RegisterGraphRoutes 注册图数据库路由
func RegisterGraphRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	
	// 获取初始图谱数据
	mux.HandleFunc("GET /api/init-graph", func(w http.ResponseWriter, r *http.Request) {
		// 直接使用corsMiddleware函数（从forum.go中）
		corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
			// 获取图处理器
			graphHandler := GetGraphHandler()
			if graphHandler == nil {
				// 创建新的处理器
				graphHandler, _ = NewGraphHandler("bolt://localhost:7687", "neo4j", "hukaile5206")
			}
			
			graphData, err := graphHandler.GetInitialGraph()
			if err != nil {
				sendForumError(w, http.StatusInternalServerError, 50000, "获取图数据失败")
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(graphData)
		})(w, r)
	})
	
	// 扩展节点
	mux.HandleFunc("GET /api/expand-node/{nodeName}", func(w http.ResponseWriter, r *http.Request) {
		corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
			nodeName := r.PathValue("nodeName")
			if nodeName == "" {
				sendForumError(w, http.StatusBadRequest, 40000, "节点名称不能为空")
				return
			}
			
			graphHandler := GetGraphHandler()
			if graphHandler == nil {
				graphHandler, _ = NewGraphHandler("bolt://localhost:7687", "neo4j", "hukaile5206")
			}
			
			graphData, err := graphHandler.ExpandNode(nodeName)
			if err != nil {
				sendForumError(w, http.StatusInternalServerError, 50000, "扩展节点失败")
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(graphData)
		})(w, r)
	})
	
	return mux
}