package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"personal-blog/models"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Article{}, &models.Category{}, &models.Tag{}, &models.Comment{})
	return db
}

func TestCreateArticle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	handler := NewArticleHandler(db)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	db.Create(user)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})
	router.POST("/articles", handler.CreateArticle)

	// 测试正常创建文章
	articleData := map[string]interface{}{
		"title":   "测试文章",
		"content": "这是测试内容",
		"summary": "测试摘要",
		"tags":    []string{"测试", "文章"},
		"status":  "published",
	}

	jsonData, _ := json.Marshal(articleData)
	req := httptest.NewRequest("POST", "/articles", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "测试文章", response["title"])
}

func TestCreateArticleValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	handler := NewArticleHandler(db)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})
	router.POST("/articles", handler.CreateArticle)

	// 测试缺少标题
	invalidData := map[string]interface{}{
		"content": "这是测试内容",
		"status":  "published",
	}

	jsonData, _ := json.Marshal(invalidData)
	req := httptest.NewRequest("POST", "/articles", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"].(string), "标题不能为空")
}