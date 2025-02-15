package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mealplanner/config"
	"mealplanner/models"
	"mealplanner/routes"
)

func autoMigrate(db *gorm.DB) {
	// Auto migrate semua model yang diperlukan
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
}

func main() {
	// Muat file .env agar variabel lingkungan tersedia
	if err := godotenv.Load(); err != nil {
		log.Println("Tidak menemukan file .env, pastikan variabel lingkungan sudah diset")
	}

	// Muat konfigurasi dari environment
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Hubungkan ke PostgreSQL menggunakan GORM.
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database
	autoMigrate(db)

	// Buat instance Gin.
	router := gin.Default()

	// Setup semua route, termasuk endpoint user management dan Spoonacular.
	// Pastikan cfg.HTTPTimeout sudah di-set, misalnya 10s.
	routes.SetupRoutes(router, db, cfg.JwtSecret, cfg.SpoonacularApiKey, cfg.SpoonacularBaseUrl, cfg.HTTPTimeout)

	// Gunakan port 2001
	addr := ":2001"
	log.Printf("Server running on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
