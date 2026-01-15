package service

import (
	"ramah-disabilitas-be/internal/repository"
)

type DashboardStats struct {
	CourseCount     int64   `json:"course_count"`
	StudentCount    int64   `json:"student_count"`
	UngradedCount   int64   `json:"ungraded_count"`
	AverageProgress float64 `json:"average_progress"`
}

func GetLecturerDashboardStats(teacherID uint64) (*DashboardStats, error) {
	// 1. Get Course Count
	courseCount, err := repository.GetCourseCountByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	// 2. Get Student Count
	studentCount, err := repository.GetStudentCountByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	// 3. Get Ungraded Assignments Count
	ungradedCount, err := repository.GetUngradedAssignmentCountByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	// 4. Get Average Progress
	avgProgress, err := repository.GetAverageProgressByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		CourseCount:     courseCount,
		StudentCount:    studentCount,
		UngradedCount:   ungradedCount,
		AverageProgress: avgProgress,
	}, nil
}
