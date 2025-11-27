package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/alesio/gestion-actividades-deportivas/config"
	"github.com/alesio/gestion-actividades-deportivas/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB initializes the MySQL connection via GORM and runs auto migrations.
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Activity{}, &models.Enrollment{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	if strings.EqualFold(cfg.AppEnv, "dev") {
		if err := Seed(db); err != nil {
			return nil, fmt.Errorf("failed to seed database: %w", err)
		}
		log.Println("development seed executed")
	}

	log.Println("database connection established and migrations executed")
	return db, nil
}
