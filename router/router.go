package router

import (
	"gin-api/handler"
	"gin-api/middleware"
	"gin-api/pkg/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup 注册所有路由
func Setup(g *gin.Engine) {
	// 全局中间件
	g.Use(middleware.CORS())
	g.Use(middleware.Logger())
	g.Use(gin.Recovery())

	// 健康检查（无需认证）
	g.GET("/api/health", func(c *gin.Context) {
		resp.OK(c, gin.H{"status": "ok"})
	})

	api := g.Group("/api")
	{
		// 公开接口（无需认证）
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)

		// 需要认证的接口
		auth := api.Group("", middleware.JWTAuth())
		{
			auth.GET("/user/me", handler.GetCurrentUser)
			auth.GET("/user/:id", handler.GetUser)
			auth.PUT("/user/:id", handler.UpdateUser)
			auth.DELETE("/user/:id", handler.DeleteUser)
			auth.GET("/users", handler.ListUsers)
		}
	}

	// 404 处理
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, resp.R{Code: 404, Msg: "接口不存在"})
	})
}
