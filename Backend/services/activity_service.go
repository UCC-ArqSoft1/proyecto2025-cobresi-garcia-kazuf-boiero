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
	if err := s.populateAvailability(slicePointers(activities)...); err != nil {
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
	if err := s.populateAvailability(slicePointers(activities)...); err != nil {
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
	if err := s.populateAvailability(&activity); err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s *ActivityService) CreateActivity(activity *models.Activity) error {
	if err := s.db.Create(activity).Error; err != nil {
		return err
	}
	activity.EnrolledCount = 0
	activity.AvailableSlots = activity.Capacity
	return nil
}

func (s *ActivityService) UpdateActivity(activity *models.Activity) error {
	if err := s.db.Save(activity).Error; err != nil {
		return err
	}
	return s.populateAvailability(activity)
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

func (s *ActivityService) populateAvailability(activities ...*models.Activity) error {
	idSet := make([]uint, 0, len(activities))
	for _, activity := range activities {
		if activity == nil || activity.ID == 0 {
			continue
		}
		idSet = append(idSet, activity.ID)
	}
	if len(idSet) == 0 {
		return nil
	}

	type counter struct {
		ActivityID uint
		Count      int64
	}
	var counters []counter
	if err := s.db.Model(&models.Enrollment{}).
		Select("activity_id, COUNT(*) as count").
		Where("activity_id IN ? AND status = ?", idSet, "inscripto").
		Group("activity_id").
		Find(&counters).Error; err != nil {
		return err
	}

	enrolledMap := make(map[uint]int64, len(counters))
	for _, c := range counters {
		enrolledMap[c.ActivityID] = c.Count
	}

	for _, activity := range activities {
		if activity == nil {
			continue
		}
		enrolled := enrolledMap[activity.ID]
		available := activity.Capacity - int(enrolled)
		if available < 0 {
			available = 0
		}
		activity.EnrolledCount = int(enrolled)
		activity.AvailableSlots = available
	}
	return nil
}

func slicePointers(activities []models.Activity) []*models.Activity {
	if len(activities) == 0 {
		return nil
	}
	ptrs := make([]*models.Activity, 0, len(activities))
	for i := range activities {
		ptrs = append(ptrs, &activities[i])
	}
	return ptrs
}
