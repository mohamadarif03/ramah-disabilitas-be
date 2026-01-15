package service

import (
	"errors"
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
	"time"
)

type AssignmentInput struct {
	Title       string    `json:"title" form:"title" binding:"required"`
	Instruction string    `json:"instruction" form:"instruction" binding:"required"`
	ModuleID    *uint64   `json:"module_id" form:"module_id"`
	MaxPoints   int       `json:"max_points" form:"max_points" binding:"required"`
	Deadline    time.Time `json:"deadline" form:"deadline" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	AllowFile   bool      `json:"allow_file" form:"allow_file"`
	AllowText   bool      `json:"allow_text" form:"allow_text"`
	AllowLate   bool      `json:"allow_late" form:"allow_late"`
}

func CreateAssignment(courseID uint64, input AssignmentInput, teacherID uint64) (*model.Assignment, error) {
	// Verify Course Ownership
	course, err := repository.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke kelas ini")
	}

	// Verify Module if provided
	if input.ModuleID != nil {
		module, err := repository.GetModuleByID(*input.ModuleID)
		if err != nil {
			return nil, errors.New("modul tidak ditemukan")
		}
		if module.CourseID != courseID {
			return nil, errors.New("modul tidak valid untuk kelas ini")
		}
	}

	assignment := &model.Assignment{
		CourseID:    courseID,
		ModuleID:    input.ModuleID,
		Title:       input.Title,
		Instruction: input.Instruction,
		MaxPoints:   input.MaxPoints,
		Deadline:    input.Deadline,
		AllowFile:   input.AllowFile,
		AllowText:   input.AllowText,
		AllowLate:   input.AllowLate,
		// Defaulting AllowVoice based on text/file logic or keeping it false for now as not in UI
		AllowVoice: false,
	}

	if err := repository.CreateAssignment(assignment); err != nil {
		return nil, err
	}

	return assignment, nil
}

func GetAssignmentsByCourse(courseID uint64, teacherID uint64) ([]model.Assignment, error) {
	course, err := repository.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke kelas ini")
	}

	return repository.GetAssignmentsByCourseID(courseID)
}

func GetStudentAssignments(studentID uint64, statusFilter string) ([]model.Assignment, error) {
	return repository.GetAssignmentsByStudentID(studentID, statusFilter)
}

func GetStudentAssignmentsByCourse(courseID uint64, studentID uint64) ([]model.Assignment, error) {
	// Verify student is enrolled
	inCourse, err := repository.IsStudentInCourse(courseID, studentID)
	if err != nil {
		return nil, err
	}
	if !inCourse {
		return nil, errors.New("unauthorized: anda belum bergabung di kelas ini")
	}

	return repository.GetAssignmentsByCourseID(courseID)
}

func GetAssignmentDetail(assignmentID, userID uint64) (*model.Assignment, error) {
	assignment, err := repository.GetAssignmentByID(assignmentID)
	if err != nil {
		return nil, errors.New("tugas tidak ditemukan")
	}

	course, err := repository.GetCourseByID(assignment.CourseID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	// Check permissions
	isTeacher := course.TeacherID == userID
	isStudent, err := repository.IsStudentInCourse(course.ID, userID)
	if err != nil {
		return nil, err
	}

	if !isTeacher && !isStudent {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke tugas ini")
	}

	// If Student, hide other submissions and set MySubmission
	if !isTeacher {
		for _, s := range assignment.Submissions {
			if s.StudentID == userID {
				mySub := s
				assignment.MySubmission = &mySub
				break
			}
		}
		assignment.Submissions = nil // Clear list for student
	}

	return assignment, nil
}

type GradeInput struct {
	Grade    float64 `json:"grade"` // Remove binding required as 0 is valid. Use explicit validation if needed. But binding:"required" fails on 0 for some validators? No, usually valid. But let's be safe.
	Feedback string  `json:"feedback"`
}

func GradeSubmission(submissionID uint64, input GradeInput, teacherID uint64) (*model.Submission, error) {
	// 1. Get Submission
	submission, err := repository.GetSubmissionByID(submissionID)
	if err != nil {
		return nil, errors.New("submission tidak ditemukan")
	}

	// 2. Get Assignment & Course to verify ownership
	assignment, err := repository.GetAssignmentByID(submission.AssignmentID)
	if err != nil {
		return nil, errors.New("assignment not found")
	}

	course, err := repository.GetCourseByID(assignment.CourseID)
	if err != nil {
		return nil, errors.New("course not found")
	}

	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke kelas ini")
	}

	// 3. Update Grade
	// Validate Grade vs MaxPoints
	if input.Grade < 0 || input.Grade > float64(assignment.MaxPoints) {
		// Just simplified error message
		return nil, errors.New("nilai tidak valid (melebihi batas maksimal)")
	}

	submission.Grade = input.Grade
	submission.Feedback = input.Feedback

	if err := repository.UpdateSubmission(submission); err != nil {
		return nil, err
	}

	return submission, nil
}
