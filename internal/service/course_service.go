package service

import (
	"errors"
	"math/rand"
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
	"time"
)

type CourseInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`  // Nullable/Optional by default in Go string if not sent
	ClassCode   string `json:"class_code"` // Optional for auto-generate
}

func CreateCourse(input CourseInput, teacherID uint64) (*model.Course, error) {
	// Auto-generate ClassCode if empty
	if input.ClassCode == "" {
		input.ClassCode = generateClassCode() // Helper function needed
	}

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

func generateClassCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func JoinCourse(classCode string, studentID uint64) error {
	course, err := repository.GetCourseByClassCode(classCode)
	if err != nil {
		return errors.New("kelas tidak ditemukan")
	}

	exists, err := repository.IsStudentInCourse(course.ID, studentID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("anda sudah bergabung di kelas ini")
	}

	if course.TeacherID == studentID {
		return errors.New("anda adalah pengajar di kelas ini")
	}

	return repository.AddStudentToCourse(course.ID, studentID)
}
