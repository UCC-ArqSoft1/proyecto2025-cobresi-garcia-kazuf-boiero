package models

import "time"

// User represents gym members and admins interacting with the system.
type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Name         string    `gorm:"size:255" json:"name"`
    Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"size:255;not null" json:"-"`
    Role         string    `gorm:"type:enum('socio','admin');default:'socio';not null" json:"role"`
    Enrollments  []Enrollment
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
