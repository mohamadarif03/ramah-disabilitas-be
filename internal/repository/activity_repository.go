package repository

import (
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/pkg/database"
)

func CreateActivity(activity *model.Activity) error {
	return database.DB.Create(activity).Error
}

func GetActivitiesByTeacherID(teacherID uint64, limit int) ([]model.Activity, error) {
	var activities []model.Activity

	// Join with Course to check TeacherID
	err := database.DB.Joins("JOIN courses ON courses.id = activities.course_id").
		Where("courses.teacher_id = ?", teacherID).
		Preload("User").
		Preload("Course").
		Order("activities.created_at desc").
		Limit(limit).
		Find(&activities).Error

	return activities, err
}

func GetActivitiesByTeacherIDWithPagination(teacherID uint64, limit int, offset int) ([]model.Activity, int64, error) {
	var activities []model.Activity
	var total int64

	query := database.DB.Model(&model.Activity{}).
		Joins("JOIN courses ON courses.id = activities.course_id").
		Where("courses.teacher_id = ?", teacherID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("User").
		Preload("Course").
		Order("activities.created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error

	return activities, total, err
}

func GetActivitiesByCourseID(courseID uint64, limit int) ([]model.Activity, error) {
	var activities []model.Activity

	err := database.DB.Where("course_id = ?", courseID).
		Preload("User").
		Preload("Course").
		Order("created_at desc").
		Limit(limit).
		Find(&activities).Error

	return activities, err
}
