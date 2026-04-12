package database

import (
	"database/sql"
	"log"
	"url-shortener/models"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Инициализация БД
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./urls.db")
	if err != nil {
		log.Fatal(err)
	}

	// Создаем таблицы
	createTables()
}

func createTables() {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		original_url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		clicks INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);
	
	CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	CREATE INDEX IF NOT EXISTS idx_user_id ON urls(user_id);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
}

// Пользователи
func CreateUser(email, password, name string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	result, err := DB.Exec(
		"INSERT INTO users (email, password, name) VALUES (?, ?, ?)",
		email, string(hashedPassword), name,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.User{ID: int(id), Email: email, Name: name}, nil
}

func FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := DB.QueryRow(
		"SELECT id, email, password, name, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := DB.QueryRow(
		"SELECT id, email, name, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// URL
func CreateURL(userID int, originalURL, shortCode string) (*models.URL, error) {
	result, err := DB.Exec(
		"INSERT INTO urls (user_id, original_url, short_code) VALUES (?, ?, ?)",
		userID, originalURL, shortCode,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.URL{
		ID:          int(id),
		UserID:      userID,
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}, nil
}

func FindURLByShortCode(shortCode string) (*models.URL, error) {
	url := &models.URL{}
	err := DB.QueryRow(
		"SELECT id, user_id, original_url, short_code, clicks, created_at FROM urls WHERE short_code = ?",
		shortCode,
	).Scan(&url.ID, &url.UserID, &url.OriginalURL, &url.ShortCode, &url.Clicks, &url.CreatedAt)

	if err != nil {
		return nil, err
	}
	return url, nil
}

func GetUserURLs(userID int) ([]models.URL, error) {
	rows, err := DB.Query(
		"SELECT id, user_id, original_url, short_code, clicks, created_at FROM urls WHERE user_id = ? ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ID, &url.UserID, &url.OriginalURL, &url.ShortCode, &url.Clicks, &url.CreatedAt)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func IncrementClicks(shortCode string) error {
	_, err := DB.Exec("UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?", shortCode)
	return err
}

func DeleteURL(userID int, shortCode string) error {
	_, err := DB.Exec("DELETE FROM urls WHERE user_id = ? AND short_code = ?", userID, shortCode)
	return err
}
