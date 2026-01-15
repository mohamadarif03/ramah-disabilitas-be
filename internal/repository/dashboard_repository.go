package repository

import (
	"ramah-disabilitas-be/pkg/database"
)

func GetCourseCountByTeacherID(teacherID uint64) (int64, error) {
	var count int64
	err := database.DB.Table("courses").Where("teacher_id = ?", teacherID).Count(&count).Error
	return count, err
}

func GetStudentCountByTeacherID(teacherID uint64) (int64, error) {
	var count int64
	// Count unique students enrolled in courses taught by the teacher
	err := database.DB.Table("course_students").
		Joins("JOIN courses ON course_students.course_id = courses.id").
		Where("courses.teacher_id = ?", teacherID).
		Distinct("course_students.user_id").
		Count(&count).Error
	return count, err
}

func GetUngradedAssignmentCountByTeacherID(teacherID uint64) (int64, error) {
	var count int64
	// Count submissions where grade is 0 (assuming 0 means ungraded) for courses taught by the teacher
	err := database.DB.Table("submissions").
		Joins("JOIN assignments ON submissions.assignment_id = assignments.id").
		Joins("JOIN courses ON assignments.course_id = courses.id").
		Where("courses.teacher_id = ? AND submissions.grade = 0", teacherID).
		Count(&count).Error
	return count, err
}

func GetAverageProgressByTeacherID(teacherID uint64) (float64, error) {
	// 1. Calculate Total Possible Completions (Sum of (Students in Course * Items in Course))
	type CourseStats struct {
		ID           uint64
		ItemCount    int64
		StudentCount int64
	}

	var courses []CourseStats

	// Get logic for item count per course and student count per course
	// This query gets ID, ItemCount (Materials + Assignments), and StudentCount for each course
	err := database.DB.Table("courses").
		Select(`
			courses.id, 
			(
				(SELECT COUNT(*) FROM materials m JOIN modules mod ON m.module_id = mod.id WHERE mod.course_id = courses.id) + 
				(SELECT COUNT(*) FROM assignments a WHERE a.course_id = courses.id)
			) as item_count,
			(SELECT COUNT(*) FROM course_students cs WHERE cs.course_id = courses.id) as student_count
		`).
		Where("courses.teacher_id = ?", teacherID).
		Scan(&courses).Error

	if err != nil {
		return 0, err
	}

	var totalPossible int64
	for _, c := range courses {
		totalPossible += c.ItemCount * c.StudentCount
	}

	if totalPossible == 0 {
		return 0, nil
	}

	// 2. Calculate Total Actual Completions (Material Completions + Submissions)
	var completedMaterials int64
	err = database.DB.Table("material_completions").
		Joins("JOIN materials ON material_completions.material_id = materials.id").
		Joins("JOIN modules ON materials.module_id = modules.id").
		Joins("JOIN courses ON modules.course_id = courses.id").
		Where("courses.teacher_id = ? AND material_completions.completed = ?", teacherID, true).
		Count(&completedMaterials).Error

	if err != nil {
		return 0, err
	}

	var submittedAssignments int64
	err = database.DB.Table("submissions").
		Joins("JOIN assignments ON submissions.assignment_id = assignments.id").
		Joins("JOIN courses ON assignments.course_id = courses.id").
		Where("courses.teacher_id = ?", teacherID).
		Count(&submittedAssignments).Error

	if err != nil {
		return 0, err
	}

	totalActual := completedMaterials + submittedAssignments

	// 3. Calculate Average
	average := (float64(totalActual) / float64(totalPossible)) * 100
	return average, nil
}
