package middleware

import (
	"gin-api/pkg/jwtutil"
	"gin-api/pkg/resp"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件（强制）
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp.Fail(c, 401, "请求未携带 token")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			resp.Fail(c, 401, "token 格式错误")
			c.Abort()
			return
		}

		claims, err := jwtutil.ParseToken(parts[1])
		if err != nil {
			resp.Fail(c, 401, err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) uint {
	id, _ := c.Get("user_id")
	if id == nil {
		return 0
	}
	return id.(uint)
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	name, _ := c.Get("username")
	if name == nil {
		return ""
	}
	return name.(string)
}
