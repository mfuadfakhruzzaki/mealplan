package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mealplanner/controllers"
	"mealplanner/middleware"
	"mealplanner/spoonacular"
	"time"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, jwtSecret, spoonacularApiKey, spoonacularBaseUrl string, spoonacularTimeout time.Duration) {
	// Konfigurasi CORS untuk domain tertentu.
	corsConfig := cors.Config{
		AllowOrigins:     []string{"https://eatgorithm.fuadfakhruz.id"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Tambahkan endpoint /health untuk memeriksa status aplikasi.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	authController := controllers.AuthController{DB: db, JwtSecret: jwtSecret}
	userController := controllers.UserController{DB: db}

	// Buat client dan controller untuk Spoonacular.
	spoonacularClient := spoonacular.NewClient(spoonacularApiKey, spoonacularBaseUrl, spoonacularTimeout)
	spoonacularController := controllers.SpoonacularController{
		Client: spoonacularClient,
		DB:     db,
	}

	// Route otentikasi
	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)
	r.POST("/logout", middleware.AuthMiddleware(jwtSecret), authController.Logout)

	// Route user
	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		userRoutes.GET("", userController.GetProfile)
		userRoutes.PUT("", userController.UpdateProfile)
	}

	// Route untuk meal plan Spoonacular
	mealPlanRoutes := r.Group("/mealplan")
	mealPlanRoutes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		mealPlanRoutes.GET("/generate", spoonacularController.GenerateMealPlan)
	}
}
