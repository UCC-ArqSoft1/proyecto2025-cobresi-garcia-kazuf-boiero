package services

import (
	"errors"

	"github.com/alesio/gestion-actividades-deportivas/models"
	"gorm.io/gorm"
)

var (
	ErrActivityNotFound = errors.New("activity not found")
	ErrActivityInactive = errors.New("activity is not active")
	ErrAlreadyEnrolled  = errors.New("user already enrolled in this activity")
	ErrNoCapacity       = errors.New("activity has no remaining capacity")
)

// EnrollmentService exposes enrollment use cases.
type EnrollmentService interface {
	EnrollUserInActivity(userID uint, activityID uint) (*models.Enrollment, error)
	GetUserEnrollments(userID uint) ([]models.Enrollment, error)
}

type enrollmentService struct {
	db *gorm.DB
}

func NewEnrollmentService(db *gorm.DB) EnrollmentService {
	return &enrollmentService{db: db}
}

func (s *enrollmentService) EnrollUserInActivity(userID, activityID uint) (*models.Enrollment, error) {
	var activity models.Activity
	if err := s.db.First(&activity, activityID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}

	if !activity.IsActive {
		return nil, ErrActivityInactive
	}

	// Check duplicate enrollment with active status.
	var existing models.Enrollment
	if err := s.db.Where("user_id = ? AND activity_id = ? AND status = ?", userID, activityID, "inscripto").First(&existing).Error; err == nil {
		return nil, ErrAlreadyEnrolled
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Validate remaining capacity.
	var count int64
	if err := s.db.Model(&models.Enrollment{}).
		Where("activity_id = ? AND status = ?", activityID, "inscripto").
		Count(&count).Error; err != nil {
		return nil, err
	}

	if int(count) >= activity.Capacity {
		return nil, ErrNoCapacity
	}

	enrollment := models.Enrollment{UserID: userID, ActivityID: activityID, Status: "inscripto"}
	if err := s.db.Create(&enrollment).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (s *enrollmentService) GetUserEnrollments(userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := s.db.Preload("Activity").
		Where("user_id = ? AND status = ?", userID, "inscripto").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}
