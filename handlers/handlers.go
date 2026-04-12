package handlers

import (
	"net/http"
	"url-shortener/database"
	"url-shortener/middleware"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Регистрация
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.CreateUser(req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Пользователь создан",
		"token":   token,
		"user":    user,
	})
}

// Вход
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.FindUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// Создание короткой ссылки
func CreateShortURL(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req models.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерируем короткий код если не предоставлен
	shortCode := req.ShortCode
	if shortCode == "" {
		shortCode = uuid.New().String()[:8]
	}

	url, err := database.CreateURL(userID.(int), req.OriginalURL, shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания ссылки"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Ссылка сокращена",
		"short_url":    "http://localhost:8080/s/" + url.ShortCode,
		"short_code":   url.ShortCode,
		"original_url": url.OriginalURL,
	})
}

// Получение всех ссылок пользователя
func GetUserURLs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	urls, err := database.GetUserURLs(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения ссылок"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"urls": urls})
}

// Удаление ссылки
func DeleteURL(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	shortCode := c.Param("shortCode")
	err := database.DeleteURL(userID.(int), shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ссылка не найдена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ссылка удалена"})
}

// Редирект по короткой ссылке
func RedirectToOriginal(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := database.FindURLByShortCode(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ссылка не найдена"})
		return
	}

	// Увеличиваем счетчик кликов
	go database.IncrementClicks(shortCode)

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

// Информация о ссылке
func GetURLInfo(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := database.FindURLByShortCode(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ссылка не найдена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original_url": url.OriginalURL,
		"clicks":       url.Clicks,
		"created_at":   url.CreatedAt,
	})
}
