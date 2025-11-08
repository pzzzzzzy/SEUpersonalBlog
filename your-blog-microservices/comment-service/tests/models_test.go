package tests

import (
	"testing"

	"comment-service/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCommentModel(t *testing.T) {
	// 使用内存SQLite数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	
	// 自动迁移
	err = models.AutoMigrate(db)
	assert.NoError(t, err)
	
	// 测试创建评论
	comment := &models.Comment{
		ArticleID: 1,
		UserID:    1,
		Username:  "testuser",
		Content:   "这是一条测试评论",
	}
	
	err = models.CreateComment(db, comment)
	assert.NoError(t, err)
	assert.NotZero(t, comment.ID) // 确保ID已设置
	
	// 测试根据文章ID获取评论
	comments, err := models.GetCommentsByArticleID(db, 1)
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, comment.ID, comments[0].ID)
	assert.Equal(t, "这是一条测试评论", comments[0].Content)
	
	// 测试根据ID获取评论
	retrievedComment, err := models.GetCommentByID(db, comment.ID)
	assert.NoError(t, err)
	assert.Equal(t, comment.ID, retrievedComment.ID)
	assert.Equal(t, "testuser", retrievedComment.Username)
	
	// 测试获取不存在的评论
	_, err = models.GetCommentByID(db, 999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	
	// 测试删除评论
	err = models.DeleteComment(db, comment.ID)
	assert.NoError(t, err)
	
	// 验证评论已删除
	_, err = models.GetCommentByID(db, comment.ID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	
	// 测试获取用户评论
	// 创建更多评论
	comment2 := &models.Comment{
		ArticleID: 2,
		UserID:    1,
		Username:  "testuser",
		Content:   "这是用户的第二条评论",
	}
	
	comment3 := &models.Comment{
		ArticleID: 3,
		UserID:    1,
		Username:  "testuser",
		Content:   "这是用户的第三条评论",
	}
	
	err = models.CreateComment(db, comment2)
	assert.NoError(t, err)
	err = models.CreateComment(db, comment3)
	assert.NoError(t, err)
	
	// 获取用户评论（分页）
	userComments, total, err := models.GetUserComments(db, 1, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total) // 应该有2条评论
	assert.Len(t, userComments, 2)  // 返回2条评论
}
