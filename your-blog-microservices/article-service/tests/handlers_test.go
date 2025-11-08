package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"article-service/config"
	"article-service/handlers"
	"article-service/models"
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
	articleHandler := handlers.NewArticleHandler(db, testConfig)
	
	// 设置不需要认证的路由
	r.GET("/api/articles", articleHandler.GetArticles)
	r.GET("/api/articles/:id", articleHandler.GetArticle)
	
	// 需要认证的路由
	authGroup := r.Group("/")
	authGroup.Use(mockAuthMiddleware())
	authGroup.POST("/api/articles", articleHandler.CreateArticle)
	authGroup.PUT("/api/articles/:id", articleHandler.UpdateArticle)
	authGroup.DELETE("/api/articles/:id", articleHandler.DeleteArticle)
	authGroup.POST("/api/articles/:id/like", articleHandler.ToggleLike)
	
	return r
}

func TestArticleHandler(t *testing.T) {
	// 使用内存SQLite数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	
	// 自动迁移
	err = models.AutoMigrate(db)
	assert.NoError(t, err)
	
	// 设置路由
	r := setupTestRouter(db)
	
	// 测试创建文章
	createPayload := `{
		"title": "测试文章标题",
		"content": "这是测试文章的内容"
	}`
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", strings.NewReader(createPayload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// 解析响应
	var createResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createResponse)
	assert.NoError(t, err)
	assert.Contains(t, createResponse, "article")
	
	// 获取创建的文章ID
	articleMap := createResponse["article"].(map[string]interface{})
	articleID := uint(articleMap["id"].(float64))
	
	// 测试获取文章列表
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/articles", nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 解析文章列表响应
	var listResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.NoError(t, err)
	assert.Contains(t, listResponse, "articles")
	articles := listResponse["articles"].([]interface{})
	assert.Len(t, articles, 1)
	
	// 测试获取单个文章
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/articles/"+strconv.FormatUint(uint64(articleID), 10), nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 解析单个文章响应
	var getResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.NoError(t, err)
	assert.Contains(t, getResponse, "article")
	
	// 测试更新文章
	updatePayload := `{
		"title": "更新后的文章标题",
		"content": "更新后的文章内容"
	}`
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/articles/"+strconv.FormatUint(uint64(articleID), 10), strings.NewReader(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 测试点赞文章
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/articles/"+strconv.FormatUint(uint64(articleID), 10)+"/like", nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 再次点赞（取消点赞）
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/articles/"+strconv.FormatUint(uint64(articleID), 10)+"/like", nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 测试删除文章
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/articles/"+strconv.FormatUint(uint64(articleID), 10), nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 验证文章已删除
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/articles/"+strconv.FormatUint(uint64(articleID), 10), nil)
	r.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}
