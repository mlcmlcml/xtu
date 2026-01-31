package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DBConfig 数据库配置
type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DB 数据库实例
var (
	db   *sql.DB
	once sync.Once
)

// GetConfig 获取数据库配置
func GetConfig() *DBConfig {
	return &DBConfig{
		Host:            getEnvOrDefault("DB_HOST", "localhost"),
		Port:            getEnvOrDefault("DB_PORT", "3306"),
		User:            getEnvOrDefault("DB_USER", "root"),
		Password:        getEnvOrDefault("DB_PASSWORD", "219332"),
		Database:        getEnvOrDefault("DB_NAME", "cybersecurity-platform"),
		MaxOpenConns:    10,  // 对应 connectionLimit: 10
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}
}

// GetDB 获取数据库连接（单例模式）
func GetDB() (*sql.DB, error) {
	var err error
	once.Do(func() {
		config := GetConfig()
		
		// 构建连接字符串
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.Database)
		
		// 打开数据库连接
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("数据库连接失败: %v", err)
			return
		}
		
		// 设置连接池参数
		db.SetMaxOpenConns(config.MaxOpenConns)
		db.SetMaxIdleConns(config.MaxIdleConns)
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
		
		// 测试连接
		if err = db.Ping(); err != nil {
			log.Printf("数据库连接测试失败: %v", err)
			return
		}
		
		log.Println("MySQL数据库连接成功！")
	})
	
	return db, err
}

// getEnvOrDefault 获取环境变量或返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		} else {
			log.Println("数据库连接已关闭")
		}
	}
}

// TestConnection 测试数据库连接
func TestConnection() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	
	// 执行简单的查询测试
	rows, err := db.Query("SELECT 1")
	if err != nil {
		return fmt.Errorf("数据库查询测试失败: %v", err)
	}
	defer rows.Close()
	
	log.Println("数据库连接测试成功！")
	return nil
}