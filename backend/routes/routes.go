package routes

import (
	"personal-blog/handlers"
	"personal-blog/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// 初始化处理器
	authHandler := handlers.NewAuthHandler(db)
	articleHandler := handlers.NewArticleHandler(db)
	commentHandler := handlers.NewCommentHandler(db)
	uploadHandler := handlers.NewUploadHandler()

	// 静态文件
	router.Static("/static", "../frontend/static")
	router.LoadHTMLGlob("../frontend/templates/*")

	// 前端页面路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// API路由
	api := router.Group("/api")
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 文章相关
		articles := api.Group("/articles")
		{
			articles.GET("", articleHandler.GetArticles)
			articles.GET("/:id", articleHandler.GetArticle)
			articles.Use(middleware.AuthMiddleware())
			{
				articles.POST("", articleHandler.CreateArticle)
				articles.PUT("/:id", articleHandler.UpdateArticle)
				articles.DELETE("/:id", articleHandler.DeleteArticle)
				articles.POST("/:id/like", articleHandler.ToggleLike)
			}
		}

		// 评论相关
		comments := api.Group("/comments")
		{
			comments.GET("/article/:articleId", commentHandler.GetComments)
			comments.Use(middleware.AuthMiddleware())
			{
				comments.POST("", commentHandler.CreateComment)
				comments.POST("/:id/like", commentHandler.LikeComment)
			}
		}

		// 文件上传
		upload := api.Group("/upload")
		{
			upload.Use(middleware.AuthMiddleware())
			upload.POST("/avatar", uploadHandler.UploadAvatar)
		}
	}
}