package model

import "time"

type ActivityType string

const (
	ActivityTypeAssignment ActivityType = "ASSIGNMENT_SUBMISSION"
	ActivityTypeMaterial   ActivityType = "MATERIAL_COMPLETION"
)

type Activity struct {
	ID          uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint64       `json:"user_id"`
	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CourseID    uint64       `json:"course_id"`
	Course      *Course      `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Type        ActivityType `json:"type"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	RelatedID   uint64       `json:"related_id"`
	CreatedAt   time.Time    `json:"created_at"`
}
