package services

import (
    "errors"

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

// AdminActivityFilter extends ActivityFilter to allow filtering by status.
type AdminActivityFilter struct {
    ActivityFilter
    IsActive *bool
}

func (s *ActivityService) ListActivities(filter ActivityFilter) ([]models.Activity, error) {
    query := s.db.Model(&models.Activity{}).Where("is_active = ?", true)
    query = applyActivityFilters(query, filter)

    var activities []models.Activity
    if err := query.Find(&activities).Error; err != nil {
        return nil, err
    }
    return activities, nil
}

func (s *ActivityService) ListActivitiesAdmin(filter AdminActivityFilter) ([]models.Activity, error) {
    query := s.db.Model(&models.Activity{})
    if filter.IsActive != nil {
        query = query.Where("is_active = ?", *filter.IsActive)
    }
    query = applyActivityFilters(query, filter.ActivityFilter)

    var activities []models.Activity
    if err := query.Find(&activities).Error; err != nil {
        return nil, err
    }
    return activities, nil
}

func (s *ActivityService) GetActivityByID(id uint) (*models.Activity, error) {
    var activity models.Activity
    if err := s.db.First(&activity, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrActivityNotFound
        }
        return nil, err
    }
    return &activity, nil
}

func (s *ActivityService) CreateActivity(activity *models.Activity) error {
    return s.db.Create(activity).Error
}

func (s *ActivityService) UpdateActivity(activity *models.Activity) error {
    return s.db.Save(activity).Error
}

func (s *ActivityService) DeleteActivity(id uint) error {
    result := s.db.Model(&models.Activity{}).Where("id = ?", id).Update("is_active", false)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrActivityNotFound
    }
    return nil
}

func applyActivityFilters(query *gorm.DB, filter ActivityFilter) *gorm.DB {
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
    return query
}
