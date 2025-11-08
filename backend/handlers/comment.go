package handlers

import (
	"net/http"
	"personal-blog/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

type CreateCommentRequest struct {
	Content   string `json:"content" binding:"required"`
	ArticleID uint   `json:"article_id" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID := c.GetUint("userID")
	
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &models.Comment{
		Content:   req.Content,
		ArticleID: req.ArticleID,
		UserID:    userID,
		ParentID:  req.ParentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.Create(comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 更新文章评论数
	h.db.Model(&models.Article{}).Where("id = ?", req.ArticleID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var comments []models.Comment
	if err := h.db.Preload("User").Preload("Replies.User").Where("article_id = ? AND parent_id IS NULL", articleID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) LikeComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := h.db.Model(&models.Comment{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment liked successfully"})
}