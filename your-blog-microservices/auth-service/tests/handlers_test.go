package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"auth-service/config"
	"auth-service/handlers"
	"auth-service/middleware"
	"auth-service/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(db *gorm.DB) *gin.Engine {
	// 设置Gin模式为测试模式
	gin.SetMode(gin.TestMode)
	
	// 创建配置
	testConfig := &config.Config{
		JWTSecret: "test-secret-key",
	}
	
	// 创建路由
	r := gin.Default()
	
	// 创建处理器
	authHandler := handlers.NewAuthHandler(db, testConfig)
	
	// 设置路由
	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)
	
	// 需要认证的路由
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware(testConfig))
	authGroup.GET("/api/auth/me", authHandler.GetUserInfo)
	
	return r
}

func TestAuthHandler(t *testing.T) {
	// 使用内存SQLite数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	
	// 自动迁移
	err = models.AutoMigrate(db)
	assert.NoError(t, err)
	
	// 设置路由
	r := setupTestRouter(db)
	
	// 测试注册
	registerPayload := `{
		"username": "testuser",
		"email": "test@example.com",
		"password": "password123"
	}`
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(registerPayload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// 解析响应
	var registerResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &registerResponse)
	assert.NoError(t, err)
	assert.Contains(t, registerResponse, "user")
	assert.Contains(t, registerResponse, "token")
	
	// 测试登录
	loginPayload := `{
		"username": "testuser",
		"password": "password123"
	}`
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 解析登录响应获取token
	var loginResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.NoError(t, err)
	assert.Contains(t, loginResponse, "user")
	assert.Contains(t, loginResponse, "token")
	token := loginResponse["token"].(string)
	
	// 测试获取用户信息（需要认证）
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 解析用户信息响应
	var userInfoResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &userInfoResponse)
	assert.NoError(t, err)
	assert.Contains(t, userInfoResponse, "id")
	assert.Contains(t, userInfoResponse, "username")
	assert.Contains(t, userInfoResponse, "email")
	assert.Equal(t, "testuser", userInfoResponse["username"])
	assert.Equal(t, "test@example.com", userInfoResponse["email"])
	assert.NotContains(t, userInfoResponse, "password")
	
	// 测试错误的密码登录
	wrongPasswordPayload := `{
		"username": "testuser",
		"password": "wrongpassword"
	}`
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", strings.NewReader(wrongPasswordPayload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
