package handlers

import (
	"net/http"
	"personal-blog/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	db *gorm.DB
}

func NewArticleHandler(db *gorm.DB) *ArticleHandler {
	return &ArticleHandler{db: db}
}

type CreateArticleRequest struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Summary    string   `json:"summary"`
	CategoryID uint     `json:"category_id"`
	Tags       []string `json:"tags"`
	Status     string   `json:"status"` // draft, published
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数绑定失败: " + err.Error()})
		return
	}

	// 添加详细的参数验证
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}

	article := &models.Article{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		AuthorID:    userID,
		CategoryID:  req.CategoryID,
		Status:      req.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.Status == "" {
		article.Status = "draft"
	}

	if err := h.db.Create(article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	// 处理标签
	if len(req.Tags) > 0 {
		var tags []models.Tag
		for _, tagName := range req.Tags {
			var tag models.Tag
			h.db.FirstOrCreate(&tag, models.Tag{Name: tagName})
			tags = append(tags, tag)
		}
		h.db.Model(article).Association("Tags").Append(tags)
	}

	c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	var articles []models.Article
	query := h.db.Preload("Author").Preload("Category").Preload("Tags")

	// 获取查询参数
	status := c.Query("status")
	categoryID := c.Query("category_id")
	authorID := c.Query("author_id")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if authorID != "" {
		query = query.Where("author_id = ?", authorID)
	}

	if err := query.Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var article models.Article
	if err := h.db.Preload("Author").Preload("Category").Preload("Tags").Preload("Comments.User").First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// 增加阅读量
	h.db.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var article models.Article
	if err := h.db.Where("id = ? AND author_id = ?", id, userID).First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found or unauthorized"})
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article.Title = req.Title
	article.Content = req.Content
	article.Summary = req.Summary
	article.CategoryID = req.CategoryID
	article.Status = req.Status
	article.UpdatedAt = time.Now()

	if err := h.db.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	if err := h.db.Where("id = ? AND author_id = ?", id, userID).Delete(&models.Article{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}