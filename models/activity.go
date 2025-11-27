package models

import "time"

// Activity describes sports activities offered by the gym.
type Activity struct {
    ID          uint         `gorm:"primaryKey" json:"id"`
    Title       string       `gorm:"size:255;not null" json:"title"`
    Description string       `gorm:"type:text" json:"description"`
    Category    string       `gorm:"size:100" json:"category"`
    DayOfWeek   int          `gorm:"type:tinyint" json:"day_of_week"`
    StartTime   time.Time    `gorm:"type:time" json:"start_time"`
    EndTime     time.Time    `gorm:"type:time" json:"end_time"`
    Capacity    int          `json:"capacity"`
    Instructor  string       `gorm:"size:255" json:"instructor"`
    ImageURL    *string      `gorm:"size:512" json:"image_url"`
    IsActive    bool         `gorm:"default:true" json:"is_active"`
    Enrollments []Enrollment `json:"-"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}
