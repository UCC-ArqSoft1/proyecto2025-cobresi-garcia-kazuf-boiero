package services

import (
    "fmt"

    "github.com/alesio/gestion-actividades-deportivas/models"
    "gorm.io/gorm"
)

// EnrollmentService manages enrollment workflows.
type EnrollmentService struct {
    db *gorm.DB
}

func NewEnrollmentService(db *gorm.DB) *EnrollmentService {
    return &EnrollmentService{db: db}
}

func (s *EnrollmentService) EnrollUserInActivity(userID, activityID uint) (*models.Enrollment, error) {
    enrollment := models.Enrollment{UserID: userID, ActivityID: activityID, Status: "inscripto"}

    // TODO: Validate capacity and duplicate enrollment before creating the record.

    if err := s.db.Create(&enrollment).Error; err != nil {
        return nil, err
    }
    return &enrollment, nil
}

func (s *EnrollmentService) GetEnrollmentsByUser(userID uint) ([]models.Enrollment, error) {
    var enrollments []models.Enrollment
    if err := s.db.Preload("Activity").Where("user_id = ?", userID).Find(&enrollments).Error; err != nil {
        return nil, err
    }
    return enrollments, nil
}

func (s *EnrollmentService) CancelEnrollment(enrollmentID uint) error {
    // TODO: Business logic for canceling enrollment and adjusting capacity.
    return fmt.Errorf("cancel enrollment not implemented")
}
