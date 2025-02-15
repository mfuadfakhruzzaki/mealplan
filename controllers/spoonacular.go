package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mealplanner/models"
	"mealplanner/spoonacular"
)

// Fungsi utilitas perhitungan (Anda bisa juga memindahkannya ke package util)
func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

func getActivityFactor(activityLevel string) float64 {
	switch activityLevel {
	case "sedentary":
		return 1.2
	case "lightly active":
		return 1.375
	case "moderately active":
		return 1.55
	case "very active":
		return 1.725
	case "extra active":
		return 1.9
	default:
		return 1.2
	}
}

func calculateDailyCalorie(user models.User, age int) float64 {
	var bmr float64
	if user.Gender == "male" {
		bmr = 88.362 + (13.397 * user.Weight) + (4.799 * user.Height) - (5.677 * float64(age))
	} else if user.Gender == "female" {
		bmr = 447.593 + (9.247 * user.Weight) + (3.098 * user.Height) - (4.330 * float64(age))
	}
	return bmr * getActivityFactor(user.ActivityLevel)
}

// SpoonacularController sekarang juga menerima DB untuk mengakses data user.
type SpoonacularController struct {
	Client *spoonacular.Client
	DB     *gorm.DB
}

// GenerateMealPlan mengambil data user untuk menghitung target kalori
// secara otomatis sebelum memanggil Spoonacular API.
func (sc *SpoonacularController) GenerateMealPlan(c *gin.Context) {
	// Ambil username dari context (dimasukkan oleh middleware otentikasi)
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Ambil data user dari DB
	var user models.User
	if err := sc.DB.Where("username = ?", username.(string)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Hitung umur dan kebutuhan kalori
	age := calculateAge(user.BirthDate)
	targetCalories := calculateDailyCalorie(user, age)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Panggil Spoonacular API dengan targetCalories yang dihitung
	mealPlan, err := sc.Client.GenerateMealPlan(ctx, int(targetCalories), "day")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Untuk tiap meal, ambil detail resep termasuk nutrisi
	detailedMeals := make([]map[string]interface{}, 0)
	for _, meal := range mealPlan.Meals {
		recipe, err := sc.Client.GetRecipeInformation(ctx, meal.ID, true)
		if err != nil {
			// Lewati jika terjadi error
			continue
		}
		detailedMeal := map[string]interface{}{
			"meal":   meal,
			"recipe": recipe,
		}
		detailedMeals = append(detailedMeals, detailedMeal)
	}

	c.JSON(http.StatusOK, gin.H{
		"target_calories": targetCalories,
		"meals":           detailedMeals,
		"nutrients":       mealPlan.Nutrients,
	})
}
