# 🛡️ XTU Cybersecurity Platform
> 基于 Vue + Go 的前后端分离网络安全教育与管理平台
---
## 📖 项目简介
XTU Cybersecurity Platform 是一个前后端分离架构的网络安全平台项目，主要用于实现网络安全知识展示、用户管理及系统功能扩展。
该项目旨在锻炼前后端分离开发能力，熟悉现代 Web 架构与接口设计模式。
---
## 🏗️ 项目结构
xtu-platform
├── frontend # 前端项目（Vue）
└── backend # 后端项目（Go）
---
## 🚀 技术栈
### 前端
- Vue
- Vue Router
- Vuex
- Axios
- Node.js
- npm
### 后端
- Go
- Gin（如使用）
- RESTful API
---
## ⚙️ 前端启动方式
进入 frontend 目录：
```bash
cd frontend
npm install
npm run serve
默认运行地址：
http://localhost:8080
🔧 后端启动方式
进入 backend 目录：
cd backend
go run main.go
默认运行端口示例：
http://localhost:8081
📡 前后端交互说明
前端通过 Axios 向后端 RESTful API 发送请求，实现数据交互。
接口示例：
GET  /api/user
POST /api/login
📦 功能模块（示例）
用户登录与注册
权限管理
数据展示模块
网络安全知识模块
后台管理系统
📈 后续优化方向
JWT 鉴权机制
MySQL 数据库接入
Docker 容器化部署
Nginx 反向代理
云服务器部署
CI/CD 自动化部署
🎯 项目目标
熟悉前后端分离开发模式
掌握 Go 后端接口开发
理解 RESTful API 设计思想
提升工程化与项目结构管理能力
👨‍💻 作者
ml c
湘潭大学
📄 License
This project is for learning and academic purposes only.
