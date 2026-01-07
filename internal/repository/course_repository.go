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

func GetCourseByClassCode(code string) (*model.Course, error) {
	var course model.Course
	err := database.DB.Where("class_code = ?", code).First(&course).Error
	return &course, err
}

func AddStudentToCourse(courseID, studentID uint64) error {
	course := model.Course{ID: courseID}
	student := model.User{ID: studentID}
	return database.DB.Model(&course).Association("Students").Append(&student)
}

func IsStudentInCourse(courseID, studentID uint64) (bool, error) {
	var count int64
	// Assumes standard naming convention course_students(course_id, user_id)
	// If User struct name is User, foreign key is user_id.
	err := database.DB.Table("course_students").Where("course_id = ? AND user_id = ?", courseID, studentID).Count(&count).Error
	return count > 0, err
}
