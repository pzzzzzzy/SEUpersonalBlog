package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadPath string
}

func NewUploadHandler() *UploadHandler {
	uploadPath := "../frontend/static/uploads"
	os.MkdirAll(uploadPath, 0755)
	return &UploadHandler{uploadPath: uploadPath}
}

func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 检查文件类型
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image files are allowed (jpg, jpeg, png, gif)"})
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("avatar_%d%s", time.Now().Unix(), ext)
	filepath := filepath.Join(h.uploadPath, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 返回相对路径
	avatarURL := "/static/uploads/" + filename
	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}