package tests

import (
	"testing"

	"article-service/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestArticleModel(t *testing.T) {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移
	err = models.AutoMigrate(db)
	assert.NoError(t, err)

	// 测试创建文章
	article := &models.Article{
		Title:    "测试文章标题",
		Content:  "这是测试文章的内容",
		AuthorID: 1,
		Username: "testuser",
	}

	// 测试创建文章
	err = models.CreateArticle(db, article)
	assert.NoError(t, err)
	assert.NotZero(t, article.ID)

	// 测试根据ID获取文章
	retrievedArticle, err := models.GetArticleByID(db, article.ID)
	assert.NoError(t, err)
	assert.Equal(t, "测试文章标题", retrievedArticle.Title)
	assert.Equal(t, "这是测试文章的内容", retrievedArticle.Content)
	assert.Equal(t, uint(1), retrievedArticle.AuthorID)
	assert.Equal(t, "testuser", retrievedArticle.Username)

	// 测试获取文章列表
	articles, total, err := models.GetArticles(db, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, articles, 1)

	// 测试更新文章
	retrievedArticle.Title = "更新后的文章标题"
	retrievedArticle.Content = "更新后的文章内容"
	err = models.UpdateArticle(db, retrievedArticle)
	assert.NoError(t, err)

	// 验证更新结果
	updatedArticle, err := models.GetArticleByID(db, article.ID)
	assert.NoError(t, err)
	assert.Equal(t, "更新后的文章标题", updatedArticle.Title)
	assert.Equal(t, "更新后的文章内容", updatedArticle.Content)

	// 测试点赞功能
	// 检查初始点赞状态
	isLiked, err := models.IsLikedByUser(db, article.ID, 1)
	assert.NoError(t, err)
	assert.False(t, isLiked)

	// 测试点赞
	liked, err := models.ToggleLike(db, article.ID, 1)
	assert.NoError(t, err)
	assert.True(t, liked)

	// 验证点赞状态
	isLiked, err = models.IsLikedByUser(db, article.ID, 1)
	assert.NoError(t, err)
	assert.True(t, isLiked)

	// 测试取消点赞
	liked, err = models.ToggleLike(db, article.ID, 1)
	assert.NoError(t, err)
	assert.False(t, liked)

	// 验证取消点赞状态
	isLiked, err = models.IsLikedByUser(db, article.ID, 1)
	assert.NoError(t, err)
	assert.False(t, isLiked)

	// 测试删除文章
	err = models.DeleteArticle(db, article.ID)
	assert.NoError(t, err)

	// 验证删除结果
	_, err = models.GetArticleByID(db, article.ID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
