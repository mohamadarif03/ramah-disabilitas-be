package service

import (
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
	"time"
)

type ActivityListResponse struct {
	Activities []ActivityItemResponse `json:"activities"`
	Pagination PaginationMeta         `json:"pagination"`
}

type ActivityItemResponse struct {
	ID          uint64    `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CourseName  string    `json:"course_name"`
	StudentName string    `json:"student_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalItems  int64 `json:"total_items"`
	Limit       int   `json:"limit"`
}

func GetLecturerActivities(teacherID uint64, page, limit int) (*ActivityListResponse, error) {
	offset := (page - 1) * limit
	activities, total, err := repository.GetActivitiesByTeacherIDWithPagination(teacherID, limit, offset)
	if err != nil {
		return nil, err
	}

	var activityResponses []ActivityItemResponse
	for _, a := range activities {
		title := a.Title
		if title == "" {
			if a.Type == model.ActivityTypeAssignment {
				title = "Pengumpulan Tugas"
			} else if a.Type == model.ActivityTypeMaterial {
				title = "Penyelesaian Materi"
			}
		}

		studentName := "Unknown Student"
		if a.User != nil {
			studentName = a.User.Name
		}

		courseName := "Unknown Course"
		if a.Course != nil {
			courseName = a.Course.Title
		}

		activityResponses = append(activityResponses, ActivityItemResponse{
			ID:          a.ID,
			Type:        string(a.Type),
			Title:       title,
			Description: a.Description,
			CourseName:  courseName,
			StudentName: studentName,
			CreatedAt:   a.CreatedAt,
		})
	}

	if activityResponses == nil {
		activityResponses = []ActivityItemResponse{}
	}

	totalPage := int((total + int64(limit) - 1) / int64(limit))

	return &ActivityListResponse{
		Activities: activityResponses,
		Pagination: PaginationMeta{
			CurrentPage: page,
			TotalPage:   totalPage,
			TotalItems:  total,
			Limit:       limit,
		},
	}, nil
}
