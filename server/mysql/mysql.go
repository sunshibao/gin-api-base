package mysqlServer

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// models 统一注册需要自动建表的模型，新增模型在此追加
var models = []interface{}{}

// RegisterModel 注册需要自动迁移的模型
func RegisterModel(m ...interface{}) {
	models = append(models, m...)
}

// InitMySQL 初始化 MySQL 连接并自动建表
func InitMySQL() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.database"),
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取 MySQL 连接池失败: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("MySQL 连接成功")

	// 自动建表
	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
		log.Println("数据库迁移完成")
	}
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}
