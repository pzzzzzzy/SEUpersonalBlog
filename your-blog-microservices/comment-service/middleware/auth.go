package middleware

import (
	"comment-service/config"
	"comment-service/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 创建JWT认证中间件
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证令牌"})
			c.Abort()
			return
		}

		// 检查格式并提取令牌
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式无效"})
			c.Abort()
			return
		}

		token := parts[1]

		// 调用认证服务验证令牌
		userInfo, err := utils.GetUserFromAuthService(token, cfg.AuthService)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或过期的令牌: " + err.Error()})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", userInfo.User.ID)
		c.Set("username", userInfo.User.Username)
		c.Next()
	}
}