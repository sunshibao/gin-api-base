package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gin-api/config"
	"gin-api/model"
	"gin-api/router"
	"gin-api/server"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 1. 加载配置
	config.Init()

	// 2. 设置运行模式
	gin.SetMode(viper.GetString("runmode"))

	// 3. 初始化数据库
	server.InitMySQL()
	server.InitRedis()

	// 4. 自动建表
	if err := model.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 5. 创建路由
	g := gin.New()
	router.Setup(g)

	// 6. 启动服务
	addr := viper.GetString("addr")
	srv := &http.Server{
		Addr:           addr,
		Handler:        g,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("🚀 服务启动: http://localhost%s\n", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
