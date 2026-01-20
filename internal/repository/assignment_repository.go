package repository

import (
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/pkg/database"
	"time"
)

func CreateAssignment(assignment *model.Assignment) error {
	return database.DB.Create(assignment).Error
}

func CreateSubmission(submission *model.Submission) error {
	return database.DB.Create(submission).Error
}

func GetAssignmentsByCourseID(courseID uint64) ([]model.Assignment, error) {
	var assignments []model.Assignment
	err := database.DB.Where("course_id = ?", courseID).Find(&assignments).Error
	return assignments, err
}

func GetAssignmentByID(id uint64) (*model.Assignment, error) {
	var assignment model.Assignment
	err := database.DB.Preload("Submissions").First(&assignment, id).Error
	return &assignment, err
}

func GetAssignmentsByStudentID(studentID uint64, statusFilter string) ([]model.Assignment, error) {
	var assignments []model.Assignment
	// Join with course_students to find courses this student joined
	// Then get assignments from those courses
	query := database.DB.Table("assignments").
		Joins("JOIN courses ON assignments.course_id = courses.id").
		Joins("JOIN course_students ON courses.id = course_students.course_id").
		Where("course_students.user_id = ?", studentID)

	if statusFilter == "overdue" {
		query = query.Where("assignments.deadline < ?", time.Now())
	} else if statusFilter == "upcoming" {
		query = query.Where("assignments.deadline >= ?", time.Now())
	}

	// Sort by deadline ascending (nearest deadline first)
	err := query.Order("assignments.deadline ASC").Find(&assignments).Error
	return assignments, err
}

func GetSubmissionByStudent(assignmentID, studentID uint64) (*model.Submission, error) {
	var submission model.Submission
	err := database.DB.Where("assignment_id = ? AND student_id = ?", assignmentID, studentID).First(&submission).Error
	return &submission, err
}

func UpdateSubmission(submission *model.Submission) error {
	return database.DB.Save(submission).Error
}

func GetSubmissionByID(id uint64) (*model.Submission, error) {
	var submission model.Submission
	err := database.DB.Preload("Student").First(&submission, id).Error
	return &submission, err
}

func GetSubmissionsByAssignmentID(assignmentID uint64) ([]model.Submission, error) {
	var submissions []model.Submission
	err := database.DB.Where("assignment_id = ?", assignmentID).Preload("Student").Find(&submissions).Error
	return submissions, err
}

func UpdateAssignment(assignment *model.Assignment) error {
	return database.DB.Save(assignment).Error
}

func DeleteAssignment(id uint64) error {
	return database.DB.Delete(&model.Assignment{}, id).Error
}
