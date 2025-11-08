package tests

import (
	"testing"

	"auth-service/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserModel(t *testing.T) {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移
	err = models.AutoMigrate(db)
	assert.NoError(t, err)

	// 测试创建用户
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Avatar:   "https://example.com/avatar.jpg",
		Bio:      "Test user bio",
	}

	// 测试密码哈希
	err = user.HashPassword("password123")
	assert.NoError(t, err)

	// 测试创建用户
	err = models.CreateUser(db, user)
	assert.NoError(t, err)

	// 测试获取用户
	retrievedUser, err := models.GetUserByUsername(db, "testuser")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", retrievedUser.Username)
	assert.Equal(t, "test@example.com", retrievedUser.Email)
	assert.Equal(t, "https://example.com/avatar.jpg", retrievedUser.Avatar)
	assert.Equal(t, "Test user bio", retrievedUser.Bio)

	// 测试密码验证
	assert.True(t, retrievedUser.CheckPassword("password123"))
	assert.False(t, retrievedUser.CheckPassword("wrongpassword"))

	// 测试根据Email获取用户
	userByEmail, err := models.GetUserByEmail(db, "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", userByEmail.Username)

	// 测试根据ID获取用户
	userByID, err := models.GetUserByID(db, retrievedUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", userByID.Username)

	// 测试删除用户（软删除）
	// 这里使用GORM的Delete方法进行软删除
	result := db.Delete(&models.User{}, retrievedUser.ID)
	assert.NoError(t, result.Error)

	// 尝试获取已删除的用户（应该返回错误）
	_, err = models.GetUserByID(db, retrievedUser.ID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
