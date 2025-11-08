package handlers

import (
	"article-service/config"
	"article-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	db     *gorm.DB
	config *config.Config
}

func NewArticleHandler(db *gorm.DB, cfg *config.Config) *ArticleHandler {
	return &ArticleHandler{
		db:     db,
		config: cfg,
	}
}

// CreateArticleRequest 创建文章请求结构
type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,min=1"`
}

// UpdateArticleRequest 更新文章请求结构
type UpdateArticleRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,min=1"`
}

// CreateArticle 创建文章
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效: " + err.Error()})
		return
	}

	// 从上下文中获取用户信息
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	// 创建文章
	article := &models.Article{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userID.(uint),
		Username: username.(string),
	}

	// 保存到数据库
	if err := models.CreateArticle(h.db, article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"article": article,
	})
}

// GetArticles 获取文章列表
func (h *ArticleHandler) GetArticles(c *gin.Context) {
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

	// 获取文章列表
	articles, total, err := models.GetArticles(h.db, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	// 检查当前用户是否已点赞文章（如果已登录）
	userID, exists := c.Get("userID")
	if exists {
		for i := range articles {
			_, _ = models.IsLikedByUser(h.db, articles[i].ID, userID.(uint))
			// 在响应中添加点赞状态（这里我们不修改models，可以在返回时处理）
			// 注意：为了简化，这里实际上需要创建一个新的结构体来包含点赞状态
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":  articles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetArticle 获取文章详情
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取文章
	article, err := models.GetArticleByID(h.db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查当前用户是否已点赞文章（如果已登录）
	response := gin.H{
		"article": article,
	}

	userID, exists := c.Get("userID")
	if exists {
		liked, _ := models.IsLikedByUser(h.db, article.ID, userID.(uint))
		response["liked"] = liked
	}

	c.JSON(http.StatusOK, response)
}

// UpdateArticle 更新文章
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取文章
	article, err := models.GetArticleByID(h.db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限（只有作者可以更新文章）
	userID, _ := c.Get("userID")
	if article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改此文章"})
		return
	}

	// 解析请求参数
	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效: " + err.Error()})
		return
	}

	// 更新文章
	article.Title = req.Title
	article.Content = req.Content

	if err := models.UpdateArticle(h.db, article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"article": article,
	})
}

// DeleteArticle 删除文章
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取文章
	article, err := models.GetArticleByID(h.db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限（只有作者可以删除文章）
	userID, _ := c.Get("userID")
	if article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此文章"})
		return
	}

	// 删除文章
	if err := models.DeleteArticle(h.db, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}

// ToggleLike 切换文章点赞状态
func (h *ArticleHandler) ToggleLike(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取用户ID
	userID, _ := c.Get("userID")

	// 切换点赞状态
	liked, err := models.ToggleLike(h.db, uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	// 获取最新的文章信息
	article, err := models.GetArticleByID(h.db, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章信息失败"})
		return
	}

	message := "取消点赞成功"
	if liked {
		message = "点赞成功"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    message,
		"liked":      liked,
		"like_count": article.LikeCount,
	})
}
