package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"comment-service/config"
	"comment-service/handlers"
	"comment-service/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// mockAuthMiddleware 创建模拟认证中间件，直接设置用户信息
func mockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 模拟已认证的用户
		c.Set("userID", uint(1))
		c.Set("username", "testuser")
		c.Next()
	}
}

// mockOtherUserAuthMiddleware 创建模拟其他用户的认证中间件
func mockOtherUserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 模拟另一个已认证的用户
		c.Set("userID", uint(2))
		c.Set("username", "otheruser")
		c.Next()
	}
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	// 设置Gin模式为测试模式
	gin.SetMode(gin.TestMode)
	
	// 创建配置
	testConfig := &config.Config{
		AuthService: "http://auth-service:8081",
	}
	
	// 创建路由
	r := gin.Default()
	
	// 创建处理器
	commentHandler := handlers.NewCommentHandler(db, testConfig)
	
	// 设置不需要认证的路由
	r.GET("/api/comments/article/:article_id", commentHandler.GetCommentsByArticleID)
	
	// 需要认证的路由
	authGroup := r.Group("/")
	authGroup.Use(mockAuthMiddleware())
	authGroup.POST("/api/comments", commentHandler.CreateComment)
	authGroup.DELETE("/api/comments/:id", commentHandler.DeleteComment)
	authGroup.GET("/api/comments/user", commentHandler.GetUserComments)
	
	// 用于测试权限的路由组
	otherUserGroup := r.Group("/")
	otherUserGroup.Use(mockOtherUserAuthMiddleware())
	otherUserGroup.DELETE("/api/comments/other/:id", commentHandler.DeleteComment)
	
	return r
}

func TestCommentHandler(t *testing.T) {
	// 使用内存SQLite数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	
	// 自动迁移
	err = models.AutoMigrate(db)
	assert.NoError(t, err)
	
	// 设置路由
	r := setupTestRouter(db)
	
	// 测试创建评论
	createPayload := `{
		"article_id": 1,
		"content": "这是一条测试评论"
	}`
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/comments", strings.NewReader(createPayload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// 解析响应
	var createResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createResponse)
	assert.NoError(t, err)
	assert.Contains(t, createResponse, "comment")
	
	// 获取创建的评论ID
	commentMap := createResponse["comment"].(map[string]interface{})
	commentID := uint(commentMap["id"].(float64))
	
	// 测试获取文章的评论列表
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/comments/article/1", nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 解析评论列表响应
	var listResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.NoError(t, err)
	assert.Contains(t, listResponse, "comments")
	comments := listResponse["comments"].([]interface{})
	assert.Len(t, comments, 1)
	
	// 测试获取用户评论
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/comments/user", nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 解析用户评论响应
	var userCommentsResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &userCommentsResponse)
	assert.NoError(t, err)
	assert.Contains(t, userCommentsResponse, "comments")
	userComments := userCommentsResponse["comments"].([]interface{})
	assert.Len(t, userComments, 1)
	
	// 测试权限控制 - 其他用户无法删除
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/comments/other/"+strconv.FormatUint(uint64(commentID), 10), nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	// 测试删除评论（正确权限）
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/comments/"+strconv.FormatUint(uint64(commentID), 10), nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 验证评论已删除
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/comments/article/1", nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var afterDeleteResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &afterDeleteResponse)
	assert.NoError(t, err)
	afterDeleteComments := afterDeleteResponse["comments"].([]interface{})
	assert.Len(t, afterDeleteComments, 0)
	
	// 测试删除不存在的评论
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/comments/999", nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}
