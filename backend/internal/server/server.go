package server

import (
	"fmt"
	"personal-blog/config"
	"personal-blog/db"
	"personal-blog/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
	router *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		router: gin.Default(),
	}
}

func (s *Server) Start() error {
	// 初始化数据库
	database, err := db.InitDB(s.config)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	s.db = database

	// 添加CORS中间件
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// 设置路由
	routes.SetupRoutes(s.router, s.db)

	// 启动服务器
	addr := fmt.Sprintf(":%s", s.config.Port)
	fmt.Printf("Server starting on %s\n", addr)
	return s.router.Run(addr)
}