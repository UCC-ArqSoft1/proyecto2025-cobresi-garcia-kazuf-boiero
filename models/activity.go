package models

import "time"

// Activity describes sports activities offered by the gym.
type Activity struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string `gorm:"size:255;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Category    string `gorm:"size:100;not null" json:"category"`
	DayOfWeek   int    `gorm:"not null" json:"day_of_week"`
	StartTime   string `gorm:"size:8;not null" json:"start_time"`
	EndTime     string `gorm:"size:8;not null" json:"end_time"`
	Capacity    int    `gorm:"not null" json:"capacity"`
	Instructor  string `gorm:"size:255;not null" json:"instructor"`
	ImageURL    string `gorm:"size:512" json:"image_url"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	// Computed fields populated at runtime so the frontend can render cupos dinámicos.
	AvailableSlots int       `gorm:"-" json:"available_slots"`
	EnrolledCount  int       `gorm:"-" json:"enrolled_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Enrollments []Enrollment `gorm:"foreignKey:ActivityID" json:"-"`
}
