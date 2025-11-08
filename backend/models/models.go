package models

import (
	"time"
	"gorm.io/gorm"
)

type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ArticleID uint      `json:"article_id"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"article"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"unique;not null" json:"username"`
	Email        string         `gorm:"unique;not null" json:"email"`
	Password     string         `gorm:"not null" json:"-"`
	Avatar       string         `json:"avatar"`
	Bio          string         `json:"bio"`
	ThirdPartyID string         `json:"third_party_id"`
	Provider     string         `json:"provider"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Articles     []Article      `gorm:"foreignKey:AuthorID" json:"articles,omitempty"`
}

type Article struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Content     string         `gorm:"type:text" json:"content"`
	Summary     string         `json:"summary"`
	AuthorID    uint           `json:"author_id"`
	Author      User           `gorm:"foreignKey:AuthorID" json:"author"`
	CategoryID  uint           `json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Tags        []Tag          `gorm:"many2many:article_tags;" json:"tags"`
	Status      string         `gorm:"default:draft" json:"status"` // draft, published
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	CommentCount int           `gorm:"default:0" json:"comment_count"`
	LikeCount   int            `gorm:"default:0" json:"like_count"`
	IsDraft     bool           `gorm:"default:true" json:"is_draft"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Comments    []Comment      `gorm:"foreignKey:ArticleID" json:"comments,omitempty"`
	Likes       []Like         `gorm:"foreignKey:ArticleID" json:"likes,omitempty"`
}

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Desc      string         `json:"desc"`
	UserID    uint           `json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Articles  []Article      `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Tag struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"unique;not null" json:"name"`
	Articles []Article `gorm:"many2many:article_tags;" json:"articles,omitempty"`
}

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"not null" json:"content"`
	ArticleID uint           `json:"article_id"`
	Article   Article        `gorm:"foreignKey:ArticleID" json:"article"`
	UserID    uint           `json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	ParentID  *uint          `json:"parent_id"`
	Parent    *Comment       `gorm:"foreignKey:ParentID" json:"parent"`
	Replies   []Comment      `gorm:"foreignKey:ParentID" json:"replies"`
	Likes     int            `gorm:"default:0" json:"likes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}