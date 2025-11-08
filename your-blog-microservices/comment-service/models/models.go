package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ArticleID uint      `json:"article_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Username  string    `json:"username" gorm:"size:100;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Comment{})
}

// CreateComment 创建评论
func CreateComment(db *gorm.DB, comment *Comment) error {
	return db.Create(comment).Error
}

// GetCommentsByArticleID 获取文章的所有评论
func GetCommentsByArticleID(db *gorm.DB, articleID uint) ([]Comment, error) {
	var comments []Comment
	err := db.Where("article_id = ?", articleID).Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetCommentByID 根据ID获取评论
func GetCommentByID(db *gorm.DB, id uint) (*Comment, error) {
	var comment Comment
	err := db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteComment 删除评论
func DeleteComment(db *gorm.DB, id uint) error {
	return db.Delete(&Comment{}, id).Error
}

// GetUserComments 获取用户的所有评论
func GetUserComments(db *gorm.DB, userID uint, page, pageSize int) ([]Comment, int64, error) {
	var comments []Comment
	var total int64

	// 计算总数
	db.Model(&Comment{}).Where("user_id = ?", userID).Count(&total)

	// 查询评论列表（按创建时间倒序）
	offset := (page - 1) * pageSize
	err := db.Where("user_id = ?", userID).Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}