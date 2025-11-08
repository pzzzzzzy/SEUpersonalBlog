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

	// 获取当前用户ID（如果已登录）
	userID, isAuthenticated := c.Get("userID")

	// 创建一个包含所有文章的响应结构
	response := make([]map[string]interface{}, len(articles))
	for i := range articles {
		response[i] = map[string]interface{}{
			"id":            articles[i].ID,
			"title":         articles[i].Title,
			"summary":       articles[i].Summary,
			"content":       articles[i].Content,
			"author_id":     articles[i].AuthorID,
			"author":        articles[i].Author,
			"category_id":   articles[i].CategoryID,
			"category":      articles[i].Category,
			"tags":          articles[i].Tags,
			"status":        articles[i].Status,
			"view_count":    articles[i].ViewCount,
			"like_count":    articles[i].LikeCount,
			"created_at":    articles[i].CreatedAt,
			"updated_at":    articles[i].UpdatedAt,
			"liked":         false, // 默认值，将在下面覆盖
		}

		// 检查用户是否点赞
		if isAuthenticated {
			var likeCount int64
			h.db.Model(&models.Like{}).Where("user_id = ? AND article_id = ?", userID, articles[i].ID).Count(&likeCount)
			response[i]["liked"] = likeCount > 0
		}
	}

	c.JSON(http.StatusOK, response)
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

	// 获取当前用户ID（如果已登录）
	userID, isAuthenticated := c.Get("userID")

	// 创建响应结构，包含所有文章信息和点赞状态
	response := map[string]interface{}{
		"id":           article.ID,
		"title":        article.Title,
		"content":      article.Content,
		"summary":      article.Summary,
		"author_id":    article.AuthorID,
		"author":       article.Author,
		"category_id":  article.CategoryID,
		"category":     article.Category,
		"tags":         article.Tags,
		"status":       article.Status,
		"view_count":   article.ViewCount,
		"comment_count": article.CommentCount,
		"like_count":   article.LikeCount,
		"is_draft":     article.IsDraft,
		"created_at":   article.CreatedAt,
		"updated_at":   article.UpdatedAt,
		"comments":     article.Comments,
		"liked":        false, // 默认值，未点赞
	}

	// 检查用户是否点赞
	if isAuthenticated {
		var likeCount int64
		h.db.Model(&models.Like{}).Where("user_id = ? AND article_id = ?", userID, article.ID).Count(&likeCount)
		response["liked"] = likeCount > 0
	}

	c.JSON(http.StatusOK, response)
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

// ToggleLike 实现文章点赞/取消点赞功能
func (h *ArticleHandler) ToggleLike(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// 获取文章ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	// 查找文章
	var article models.Article
	if err := h.db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// 检查用户是否已点赞
	var like models.Like
	result := h.db.Where("user_id = ? AND article_id = ?", userID, id).First(&like)
	
	var isLiked bool
	var message string
	
	if result.Error != nil {
		// 用户未点赞，创建点赞记录
		if result.Error == gorm.ErrRecordNotFound {
			like = models.Like{
				UserID:    userID.(uint),
				ArticleID: uint(id),
			}
			
			// 开始事务
			tx := h.db.Begin()
			
			// 创建点赞记录
			if err := tx.Create(&like).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create like"})
				return
			}
			
			// 增加文章点赞数
			article.LikeCount++
			if err := tx.Save(&article).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
				return
			}
			
			// 提交事务
			tx.Commit()
			
			isLiked = true
			message = "Article liked successfully"
		} else {
			// 其他数据库错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check like status"})
			return
		}
	} else {
		// 用户已点赞，删除点赞记录
		// 开始事务
		tx := h.db.Begin()
		
		// 删除点赞记录
		if err := tx.Delete(&like).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like"})
			return
		}
		
		// 减少文章点赞数
		if article.LikeCount > 0 {
			article.LikeCount--
		}
		if err := tx.Save(&article).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
			return
		}
		
		// 提交事务
		tx.Commit()
		
		isLiked = false
		message = "Like removed successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         article.ID,
		"like_count": article.LikeCount,
		"liked":      isLiked,
		"message":    message,
	})
}