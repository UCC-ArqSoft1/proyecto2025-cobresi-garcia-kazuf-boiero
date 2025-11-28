package database

import (
	"log"

	"github.com/alesio/gestion-actividades-deportivas/models"
	"github.com/alesio/gestion-actividades-deportivas/security"
	"gorm.io/gorm"
)

// Seed inserts minimal development data to simplify local testing.
func Seed(db *gorm.DB) error {
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount == 0 {
		hashedPassword, err := security.HashPassword("contra123")
		if err != nil {
			return err
		}
		users := []models.User{
			{Name: "Admin", Email: "admin@example.com", PasswordHash: hashedPassword, Role: "admin"},
			{Name: "Socia Demo", Email: "socia@example.com", PasswordHash: hashedPassword, Role: "socio"},
		}
		if err := db.Create(&users).Error; err != nil {
			return err
		}
		log.Println("seed: created default users")
	}

	var activityCount int64
	if err := db.Model(&models.Activity{}).Count(&activityCount).Error; err != nil {
		return err
	}
	if activityCount == 0 {
		activities := []models.Activity{
			{
				Title:       "Yoga Sunrise",
				Description: "Clase de yoga matutina para movilidad y respiracion.",
				Category:    "yoga",
				DayOfWeek:   1,
				StartTime:   "07:30",
				EndTime:     "08:30",
				Capacity:    20,
				Instructor:  "Lucia Perez",
				ImageURL:    "",
			},
			{
				Title:       "Funcional",
				Description: "Entrenamiento de fuerza y acondicionamiento general.",
				Category:    "fuerza",
				DayOfWeek:   2,
				StartTime:   "18:00",
				EndTime:     "19:00",
				Capacity:    18,
				Instructor:  "Carlos Diaz",
				ImageURL:    "",
			},
			{
				Title:       "Spinning",
				Description: "Cardio de alta intensidad en bicicleta fija.",
				Category:    "cardio",
				DayOfWeek:   4,
				StartTime:   "19:30",
				EndTime:     "20:15",
				Capacity:    15,
				Instructor:  "Agus Flores",
				ImageURL:    "",
			},
		}
		if err := db.Create(&activities).Error; err != nil {
			return err
		}
		log.Println("seed: created sample activities")
	}

	return nil
}
