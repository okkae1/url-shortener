package main

import (
	"log"
	"url-shortener/database"
	"url-shortener/handlers"
	"url-shortener/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	database.InitDB()
	defer database.DB.Close()

	r := gin.Default()

	r.Static("/static", "./public")

	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	r.GET("/login", func(c *gin.Context) {
		c.File("./public/login.html")
	})

	r.GET("/register", func(c *gin.Context) {
		c.File("./public/register.html")
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.File("./public/dashboard.html")
	})

	api := r.Group("/api")
	{

		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		urls := api.Group("/urls")
		urls.Use(middleware.AuthMiddleware())
		{
			urls.POST("", handlers.CreateShortURL)
			urls.GET("", handlers.GetUserURLs)
			urls.DELETE("/:shortCode", handlers.DeleteURL)
		}
	}

	r.GET("/s/:shortCode", handlers.RedirectToOriginal)
	r.GET("/info/:shortCode", handlers.GetURLInfo)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("❌ Ошибка запуска сервера:", err)
	}
}
