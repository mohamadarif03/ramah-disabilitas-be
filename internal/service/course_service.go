package service

import (
	"errors"
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
)

type CourseInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	ClassCode   string `json:"class_code" binding:"required,min=4,max=20"`
}

func CreateCourse(input CourseInput, teacherID uint64) (*model.Course, error) {
	course := &model.Course{
		TeacherID:   teacherID,
		Title:       input.Title,
		Description: input.Description,
		Thumbnail:   input.Thumbnail,
		ClassCode:   input.ClassCode,
	}

	if err := repository.CreateCourse(course); err != nil {
		return nil, err
	}

	return course, nil
}

func GetCoursesByTeacher(teacherID uint64) ([]model.Course, error) {
	return repository.GetCoursesByTeacherID(teacherID)
}

func UpdateCourse(id uint64, input CourseInput, teacherID uint64) (*model.Course, error) {
	course, err := repository.GetCourseByID(id)
	if err != nil {
		return nil, err
	}

	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: you do not own this course")
	}

	course.Title = input.Title
	course.Description = input.Description
	course.Thumbnail = input.Thumbnail
	// ClassCode biasanya tidak diubah sembarangan, tapi jika perlu:
	if input.ClassCode != "" {
		course.ClassCode = input.ClassCode
	}

	if err := repository.UpdateCourse(course); err != nil {
		return nil, err
	}

	return course, nil
}

func DeleteCourse(id uint64, teacherID uint64) error {
	course, err := repository.GetCourseByID(id)
	if err != nil {
		return err
	}

	if course.TeacherID != teacherID {
		return errors.New("unauthorized: you do not own this course")
	}

	return repository.DeleteCourse(id)
}
