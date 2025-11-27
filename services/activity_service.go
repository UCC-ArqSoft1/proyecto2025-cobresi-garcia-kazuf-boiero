package services

import (
    "github.com/alesio/gestion-actividades-deportivas/models"
    "gorm.io/gorm"
)

// ActivityService encapsulates interactions with the activities table.
type ActivityService struct {
    db *gorm.DB
}

func NewActivityService(db *gorm.DB) *ActivityService {
    return &ActivityService{db: db}
}

// ActivityFilter captures optional search parameters for listing activities.
type ActivityFilter struct {
    Query    string
    Category string
    Day      *int
}

func (s *ActivityService) ListActivities(filter ActivityFilter) ([]models.Activity, error) {
    query := s.db.Model(&models.Activity{}).Where("is_active = ?", true)

    if filter.Query != "" {
        like := "%" + filter.Query + "%"
        query = query.Where("title LIKE ? OR description LIKE ?", like, like)
    }

    if filter.Category != "" {
        query = query.Where("category = ?", filter.Category)
    }

    if filter.Day != nil {
        query = query.Where("day_of_week = ?", *filter.Day)
    }

    var activities []models.Activity
    if err := query.Find(&activities).Error; err != nil {
        return nil, err
    }
    return activities, nil
}

func (s *ActivityService) GetActivityByID(id uint) (*models.Activity, error) {
    var activity models.Activity
    if err := s.db.First(&activity, id).Error; err != nil {
        return nil, err
    }
    return &activity, nil
}

func (s *ActivityService) CreateActivity(activity *models.Activity) error {
    // TODO: Apply business validations before persisting.
    return s.db.Create(activity).Error
}

func (s *ActivityService) UpdateActivity(activity *models.Activity) error {
    // TODO: Apply business validations before updating.
    return s.db.Save(activity).Error
}

func (s *ActivityService) DeleteActivity(id uint) error {
    // TODO: Decide whether to hard delete or just mark as inactive.
    return s.db.Delete(&models.Activity{}, id).Error
}
