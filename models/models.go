package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type URL struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	ShortCode   string `json:"short_code"`
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
