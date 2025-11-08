package handlers

import (
	"comment-service/config"
	"comment-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db     *gorm.DB
	config *config.Config
}

func NewCommentHandler(db *gorm.DB, cfg *config.Config) *CommentHandler {
	return &CommentHandler{
		db:     db,
		config: cfg,
	}
}

// CreateCommentRequest 创建评论请求结构
type CreateCommentRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required,min=1"`
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效: " + err.Error()})
		return
	}

	// 从上下文中获取用户信息
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	// 创建评论
	comment := &models.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID.(uint),
		Username:  username.(string),
		Content:   req.Content,
	}

	// 保存到数据库
	if err := models.CreateComment(h.db, comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": comment,
	})
}

// GetCommentsByArticleID 获取文章的所有评论
func (h *CommentHandler) GetCommentsByArticleID(c *gin.Context) {
	// 获取文章ID
	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取评论列表
	comments, err := models.GetCommentsByArticleID(h.db, uint(articleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"count":    len(comments),
	})
}

// DeleteComment 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	// 获取评论ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	// 获取评论
	comment, err := models.GetCommentByID(h.db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查权限（只有评论作者可以删除评论）
	userID, _ := c.Get("userID")
	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此评论"})
		return
	}

	// 删除评论
	if err := models.DeleteComment(h.db, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})
}

// GetUserComments 获取用户的评论列表
func (h *CommentHandler) GetUserComments(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, _ := c.Get("userID")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取用户评论列表
	comments, total, err := models.GetUserComments(h.db, userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments":  comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}