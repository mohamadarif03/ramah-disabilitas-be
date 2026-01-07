package model

import "time"

type PracticeSession struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `json:"user_id"`
	SubtestID uint64    `json:"subtest_id"`
	Score     int       `json:"score"`
	Correct   int       `json:"correct"`
	Wrong     int       `json:"wrong"`
	CreatedAt time.Time `json:"created_at"`

	Subtest Subtest `gorm:"foreignKey:SubtestID" json:"subtest"`
}

type QuestionReport struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID  uint64    `json:"question_id"`
	UserID      uint64    `json:"user_id"`
	ReportType  string    `gorm:"type:varchar(50)" json:"report_type"` // typo, ambiguity, wrong_answer, other
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, reviewed, resolved
	CreatedAt   time.Time `json:"created_at"`

	Question Question `gorm:"foreignKey:QuestionID" json:"question"`
	User     User     `gorm:"foreignKey:UserID" json:"user"`
}
