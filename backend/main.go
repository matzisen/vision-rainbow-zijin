package main

import (
	"embed"
	"fmt"
	"gotest/config"
	"gotest/internal/controllers"
	"gotest/internal/middleware"
	"gotest/pkg/ws"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS Middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

//go:embed dist/*
var content embed.FS

func main() {
	// 1. Anti-crash
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Critical Error:", err)
		}
	}()

	// 2. Setup paths
	workDir, _ := os.Getwd()
	fmt.Println(">>> Current working directory:", workDir)

	// 3. Init DB
	config.InitDB()

	// 4. Create uploads dir
	uploadDir := filepath.Join(workDir, "uploads")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// 5. Init WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// 6. Init Gin
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.Use(CORSMiddleware())
	r.StaticFS("/uploads", gin.Dir(uploadDir, true))

	// 7. Frontend Hosting
	frontendFS, err := fs.Sub(content, "dist")
	if err == nil {
		indexHtmlBytes, _ := fs.ReadFile(frontendFS, "index.html")
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api") {
				c.JSON(404, gin.H{"error": "API not found"})
				return
			}
			file, err := frontendFS.Open(strings.TrimPrefix(path, "/"))
			if err == nil {
				defer file.Close()
				stat, _ := file.Stat()
				if !stat.IsDir() {
					c.FileFromFS(path, http.FS(frontendFS))
					return
				}
			}
			c.Data(200, "text/html; charset=utf-8", indexHtmlBytes)
		})
	}

	// 9. Initialize Controllers (No Service injection for ChatController)
	chatController := &controllers.ChatController{Hub: hub}
	userController := new(controllers.UserController)
	productController := new(controllers.ProductController)
	fileController := new(controllers.FileController)
	orderController := new(controllers.OrderController)
	cartController := new(controllers.CartController)
	adminController := new(controllers.AdminController)

	// 10. API Routes
	api := r.Group("/api")
	{
		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)

		// WebSocket Endpoint
		api.GET("/ws", chatController.Connect)

		api.GET("/products", productController.List)
		api.GET("/products/:id", productController.GetDetail)
		api.GET("/categories", productController.Categories)

		// Protected Routes
		userGroup := api.Group("/", middleware.Auth())
		{
			userGroup.GET("/user/data", userController.GetMyData)
			userGroup.PUT("/user/profile", userController.UpdateProfile)
			userGroup.PUT("/user/password", userController.ChangePassword)
			userGroup.GET("/user/favorite/check", userController.CheckFavorite)
			userGroup.POST("/user/favorite", userController.ToggleFavorite)

			// ★★★ 新增路由：获取指定用户信息 (用于聊天显示) ★★★
			userGroup.GET("/users/:id", userController.GetUserInfo)

			chatGroup := userGroup.Group("/chat")
			{
				chatGroup.GET("/contacts", chatController.GetContacts)
				chatGroup.GET("/messages", chatController.GetHistory)
			}

			userGroup.POST("/upload", fileController.Upload)
			userGroup.POST("/products", productController.Create)
			userGroup.PUT("/products/:id", productController.Update)
			userGroup.POST("/orders", orderController.Create)
			userGroup.POST("/orders/batch", orderController.BatchCreate)
			userGroup.GET("/orders", orderController.List)
			userGroup.POST("/orders/:id/pay", orderController.Pay)
			userGroup.POST("/cart", cartController.Add)
			userGroup.GET("/cart", cartController.List)
			userGroup.DELETE("/cart/:id", cartController.Delete)
		}

		// Admin Routes
		adminGroup := api.Group("/admin")
		{
			adminGroup.POST("/login", adminController.Login)
			adminGroup.GET("/init", adminController.Init)
			authGroup := adminGroup.Group("/", middleware.AdminAuth())
			{
				authGroup.GET("/info", adminController.GetInfo)
				authGroup.GET("/stats", adminController.GetStats)
				authGroup.GET("/users", adminController.GetUsers)
				authGroup.PUT("/users/:id/status", adminController.UpdateUserStatus)
				authGroup.GET("/products", adminController.GetProducts)
				authGroup.PUT("/products/:id/audit", adminController.AuditProduct)
				authGroup.GET("/orders", adminController.GetOrders)
			}
		}
	}

	fmt.Println(">>> Server started successfully: http://localhost:8081")
	r.Run(":8081")
}
