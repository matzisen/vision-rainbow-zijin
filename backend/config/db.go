package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gotest/internal/models" // 确保这里的 module 名和你 go.mod 一致

	"github.com/glebarez/sqlite" // 保持使用纯 Go 驱动，避免 CGO 问题
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // ★★★ 引入日志包
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 1. 动态获取当前运行路径
	workDir, _ := os.Getwd()
	dbFileName := "xianqu.db"
	dbPath := filepath.Join(workDir, dbFileName)

	// ★★★ 2. 配置 GORM 日志 (强烈建议) ★★★
	// 这样控制台会打印出 SQL 语句，方便你调试 500 错误
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 级别：Info (打印所有 SQL)
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色打印
		},
	)

	// 3. 连接 SQLite
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger, // 挂载日志配置
	})
	if err != nil {
		log.Fatal("❌ 数据库连接失败: ", err)
	}

	// 4. 开启 SQLite 外键约束 (推荐)
	// SQLite 默认关闭外键，开启后能防止脏数据
	DB.Exec("PRAGMA foreign_keys = ON")

	// 5. 设置连接池
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	fmt.Println("✅ SQLite 数据库连接成功！路径:", dbPath)

	// 6. ★★★ 自动迁移表结构 ★★★
	// 注意：请确保 internal/models 下真的有这些结构体，否则会报错
	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},    // 订单表如果还没写，继续保持注释
		&models.Cart{},     // <--- ★★★ 这里！把前面的 // 去掉，启用购物车表
		&models.Favorite{}, // 收藏表如果写了，也可以去掉注释
		&models.Message{},
	)

	if err != nil {
		log.Fatal("❌ 自动建表失败: ", err)
	}

	fmt.Println("✅ 数据库表结构已同步完成")
}
