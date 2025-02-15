package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mealplanner/models"
)

// ProfileResponse merupakan struktur respons dengan field terhitung
type ProfileResponse struct {
	Username             string    `json:"username"`
	Email                string    `json:"email"`
	FullName             string    `json:"full_name"`
	BirthDate            time.Time `json:"birth_date"`
	Age                  int       `json:"age"`
	Gender               string    `json:"gender"`
	Weight               float64   `json:"weight"`
	Height               float64   `json:"height"`
	ActivityLevel        string    `json:"activity_level"`
	DailyCalorieRequired float64   `json:"daily_calorie_required"`
}

type UserController struct {
	DB *gorm.DB
}

func (uc *UserController) GetProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := uc.DB.Where("username = ?", username.(string)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	age := calculateAge(user.BirthDate)
	dailyCalories := calculateDailyCalorie(user, age)

	response := ProfileResponse{
		Username:             user.Username,
		Email:                user.Email,
		FullName:             user.FullName,
		BirthDate:            user.BirthDate,
		Age:                  age,
		Gender:               user.Gender,
		Weight:               user.Weight,
		Height:               user.Height,
		ActivityLevel:        user.ActivityLevel,
		DailyCalorieRequired: dailyCalories,
	}

	c.JSON(http.StatusOK, response)
}

type UpdateUserInput struct {
	Email         string  `json:"email" binding:"omitempty,email"`
	FullName      string  `json:"full_name"`
	BirthDate     string  `json:"birth_date"`     // Format: "2006-01-02"
	Gender        string  `json:"gender"`         // "male" atau "female"
	Weight        float64 `json:"weight"`         // dalam kg
	Height        float64 `json:"height"`         // dalam cm
	ActivityLevel string  `json:"activity_level"` // e.g., "sedentary", dll.
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := uc.DB.Where("username = ?", username.(string)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Email != "" {
		user.Email = input.Email
	}
	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", input.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "BirthDate harus dalam format YYYY-MM-DD"})
			return
		}
		user.BirthDate = birthDate
	}
	if input.Gender != "" {
		user.Gender = input.Gender
	}
	if input.Weight != 0 {
		user.Weight = input.Weight
	}
	if input.Height != 0 {
		user.Height = input.Height
	}
	if input.ActivityLevel != "" {
		user.ActivityLevel = input.ActivityLevel
	}

	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profil diperbarui"})
}
