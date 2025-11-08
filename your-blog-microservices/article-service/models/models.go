package models

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	AuthorID  uint      `json:"author_id" gorm:"not null"`
	Username  string    `json:"username" gorm:"size:100;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LikeCount int       `json:"like_count" gorm:"default:0"`
}

// Like 点赞模型
type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ArticleID uint      `json:"article_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Article{}, &Like{})
}

// CreateArticle 创建文章
func CreateArticle(db *gorm.DB, article *Article) error {
	return db.Create(article).Error
}

// GetArticleByID 根据ID获取文章
func GetArticleByID(db *gorm.DB, id uint) (*Article, error) {
	var article Article
	err := db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetArticles 获取文章列表（分页）
func GetArticles(db *gorm.DB, page, pageSize int) ([]Article, int64, error) {
	var articles []Article
	var total int64

	// 计算总数
	db.Model(&Article{}).Count(&total)

	// 查询文章列表（按创建时间倒序）
	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// UpdateArticle 更新文章
func UpdateArticle(db *gorm.DB, article *Article) error {
	return db.Save(article).Error
}

// DeleteArticle 删除文章
func DeleteArticle(db *gorm.DB, id uint) error {
	return db.Delete(&Article{}, id).Error
}

// IsLikedByUser 检查用户是否已点赞文章
func IsLikedByUser(db *gorm.DB, articleID, userID uint) (bool, error) {
	var count int64
	err := db.Model(&Like{}).Where("article_id = ? AND user_id = ?", articleID, userID).Count(&count).Error
	return count > 0, err
}

// ToggleLike 切换文章点赞状态
func ToggleLike(db *gorm.DB, articleID, userID uint) (bool, error) {
	var like Like

	// 检查是否已点赞
	result := db.Where("article_id = ? AND user_id = ?", articleID, userID).First(&like)

	// 如果已点赞，则取消点赞
	if result.Error == nil {
		// 开启事务
		tx := db.Begin()

		// 删除点赞记录
		if err := tx.Delete(&like).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		// 减少文章点赞数
		if err := tx.Model(&Article{}).Where("id = ?", articleID).Update("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		tx.Commit()
		return false, nil // 返回false表示已取消点赞
	}

	// 如果未点赞且错误是记录不存在，则添加点赞
	if result.Error == gorm.ErrRecordNotFound {
		// 开启事务
		tx := db.Begin()

		// 创建点赞记录
		like = Like{
			ArticleID: articleID,
			UserID:    userID,
		}
		if err := tx.Create(&like).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		// 增加文章点赞数
		if err := tx.Model(&Article{}).Where("id = ?", articleID).Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		tx.Commit()
		return true, nil // 返回true表示已点赞
	}

	// 其他错误
	return false, result.Error
}