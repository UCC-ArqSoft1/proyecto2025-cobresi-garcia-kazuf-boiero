package models

import "time"

// Enrollment links a user with an activity.
type Enrollment struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    UserID     uint      `gorm:"not null;index;uniqueIndex:idx_user_activity" json:"user_id"`
    ActivityID uint      `gorm:"not null;index;uniqueIndex:idx_user_activity" json:"activity_id"`
    Status     string    `gorm:"type:enum('inscripto','cancelado');default:'inscripto';not null" json:"status"`
    User       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
    Activity   Activity  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"activity"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

func (Enrollment) TableName() string {
    return "enrollments"
}
