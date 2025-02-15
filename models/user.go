package models

import "time"
import "gorm.io/gorm"

// User mendefinisikan model pengguna yang diperluas
type User struct {
	gorm.Model
	Username      string    `gorm:"uniqueIndex;not null" json:"username"`
	Password      string    `gorm:"not null" json:"-"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	FullName      string    `json:"full_name"`
	BirthDate     time.Time `json:"birth_date"`     // Format: YYYY-MM-DD
	Gender        string    `json:"gender"`         // "male" atau "female"
	Weight        float64   `json:"weight"`         // dalam kg
	Height        float64   `json:"height"`         // dalam cm
	ActivityLevel string    `json:"activity_level"` // e.g., "sedentary", "lightly active", "moderately active", "very active", "extra active"
}
