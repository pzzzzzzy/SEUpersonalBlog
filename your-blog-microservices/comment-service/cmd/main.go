package main

import (
	"comment-service/config"
	"comment-service/db"
	"comment-service/handlers"
	"comment-service/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移数据库表结构 - 已通过初始化脚本创建，禁用自动迁移
	// if err := models.AutoMigrate(database); err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }
	log.Println("Database tables already created via init script")

	// 初始化Consul客户端
	consulClient, err := api.NewClient(&api.Config{
		Address: cfg.ConsulAddr,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Consul: %v", err)
	}

	// 注册服务到Consul
	registration := &api.AgentServiceRegistration{
		ID:      cfg.ServiceID,
		Name:    cfg.ServiceName,
		Address: cfg.ServiceName, // 使用服务名作为地址，Docker网络中可解析
		Port:    8083,
		Tags:    []string{"comment", "api"},
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%s/health", cfg.ServiceName, cfg.ServicePort),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	if err := consulClient.Agent().ServiceRegister(registration); err != nil {
		log.Fatalf("Failed to register service to Consul: %v", err)
	}
	defer func() {
		if err := consulClient.Agent().ServiceDeregister(cfg.ServiceID); err != nil {
			log.Printf("Failed to deregister service: %v", err)
		}
	}()

	// 初始化Gin引擎
	router := gin.Default()

	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// 初始化处理器
	commentHandler := handlers.NewCommentHandler(database, cfg)

	// 设置路由
	api := router.Group("/api")
	{
		// 公开路由（无需认证）
		api.GET("/articles/:article_id/comments", commentHandler.GetCommentsByArticleID)

		// 需要认证的路由
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware(cfg))
		{
			// 评论相关
			authorized.POST("/comments", commentHandler.CreateComment)
			authorized.DELETE("/comments/:id", commentHandler.DeleteComment)
			authorized.GET("/user/comments", commentHandler.GetUserComments)
		}
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// 启动服务器（非阻塞）
	go func() {
		fmt.Printf("Comment Service starting on port %s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置5秒的超时时间来关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}