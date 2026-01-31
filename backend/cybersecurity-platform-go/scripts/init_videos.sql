-- 这个文件用来创建videos表和插入测试数据

-- 1. 创建videos表（如果不存在）
-- 注意：表结构要和Node.js项目中的一样
CREATE TABLE IF NOT EXISTS videos (
    id INT PRIMARY KEY AUTO_INCREMENT,  -- 视频ID，自动增长
    url VARCHAR(500) NOT NULL,          -- 视频URL地址
    description TEXT,                   -- 视频描述
    duration INT DEFAULT 0,             -- 视频时长（秒）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 创建时间
);

-- 2. 清空表中的旧数据（如果存在）
-- DELETE FROM videos;  -- 注意：这会删除所有已有数据

-- 3. 插入测试数据
-- 插入3条测试视频记录
INSERT INTO videos (url, description, duration) VALUES
('http://localhost:3000/api/videoing/video1.mp4', '网络安全基础教程 - 第一部分', 3600),
('http://localhost:3000/api/videoing/video2.mp4', '密码学原理与应用', 4200),
('http://localhost:3000/api/videoing/video3.mp4', '网络攻防实战演示', 5400);

-- 4. 显示插入的数据（验证）
SELECT * FROM videos;