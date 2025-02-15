package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mealplanner/models"
)

type AuthController struct {
	DB        *gorm.DB
	JwtSecret string
}

type RegisterInput struct {
	Username      string  `json:"username" binding:"required"`
	Password      string  `json:"password" binding:"required"`
	Email         string  `json:"email" binding:"required,email"`
	FullName      string  `json:"full_name"`
	BirthDate     string  `json:"birth_date" binding:"required"`     // Format: "2006-01-02"
	Gender        string  `json:"gender" binding:"required"`         // "male" atau "female"
	Weight        float64 `json:"weight" binding:"required"`         // dalam kg
	Height        float64 `json:"height" binding:"required"`         // dalam cm
	ActivityLevel string  `json:"activity_level" binding:"required"` // e.g., "sedentary", "lightly active", dll.
}

func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parsing BirthDate
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BirthDate harus dalam format YYYY-MM-DD"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}

	user := models.User{
		Username:      input.Username,
		Password:      string(hashedPassword),
		Email:         input.Email,
		FullName:      input.FullName,
		BirthDate:     birthDate,
		Gender:        input.Gender,
		Weight:        input.Weight,
		Height:        input.Height,
		ActivityLevel: input.ActivityLevel,
	}

	if err := ac.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Jangan mengembalikan password
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := ac.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ac.JwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (ac *AuthController) Logout(c *gin.Context) {
	// Logout dengan JWT umumnya dikelola di sisi klien.
	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}
