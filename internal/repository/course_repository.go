package repository

import (
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/pkg/database"
)

func CreateCourse(course *model.Course) error {
	return database.DB.Create(course).Error
}

func GetCoursesByTeacherID(teacherID uint64) ([]model.Course, error) {
	var courses []model.Course
	err := database.DB.Where("teacher_id = ?", teacherID).Find(&courses).Error
	return courses, err
}

func GetCourseByID(id uint64) (*model.Course, error) {
	var course model.Course
	err := database.DB.First(&course, id).Error
	return &course, err
}

func UpdateCourse(course *model.Course) error {
	return database.DB.Save(course).Error
}

func DeleteCourse(id uint64) error {
	return database.DB.Delete(&model.Course{}, id).Error
}
