-- 创建students表（如果不存在）
CREATE TABLE IF NOT EXISTS students (
    stuId VARCHAR(50) PRIMARY KEY,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE
);

-- 创建userdetail表（如果不存在）
CREATE TABLE IF NOT EXISTS userdetail (
    id INT PRIMARY KEY AUTO_INCREMENT,
    stuId VARCHAR(50) NOT NULL UNIQUE,
    nickName VARCHAR(50),
    userHead VARCHAR(500),
    userName VARCHAR(50),
    userEmail VARCHAR(100),
    FOREIGN KEY (stuId) REFERENCES students(stuId)
);

-- 插入测试用户数据
-- 密码都是 "123456" 经过bcrypt加密后的值
INSERT IGNORE INTO students (stuId, password, email) VALUES
('20230001', '$2a$10$N9qo8uLOickgx2ZMRZoMye3Z7c3K3K9Z7mZQ7JZkFvJX9vY7qXqZC', 'student1@example.com'),
('20230002', '$2a$10$N9qo8uLOickgx2ZMRZoMye3Z7c3K3K9Z7mZQ7JZkFvJX9vY7qXqZC', 'student2@example.com');

INSERT IGNORE INTO userdetail (stuId, nickName, userHead, userName, userEmail) VALUES
('20230001', '小明', 'https://example.com/avatar1.jpg', '张三', 'student1@example.com'),
('20230002', '小红', 'https://example.com/avatar2.jpg', '李四', 'student2@example.com');

-- 验证数据
SELECT 
    s.stuId,
    s.email,
    u.nickName,
    u.userName,
    u.userHead
FROM students s
JOIN userdetail u ON s.stuId = u.stuId;