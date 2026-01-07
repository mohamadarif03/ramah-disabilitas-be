package model

import "time"

type Match struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Player1ID uint64    `json:"player1_id"`
	Player2ID uint64    `json:"player2_id"`
	WinnerID  *uint64   `json:"winner_id"`
	SubtestID uint64    `json:"subtest_id"`
	Status    string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, ongoing, finished
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Player1 User    `gorm:"foreignKey:Player1ID" json:"player1"`
	Player2 User    `gorm:"foreignKey:Player2ID" json:"player2"`
	Subtest Subtest `gorm:"foreignKey:SubtestID" json:"subtest"`
}

type MatchDetail struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchID    uint64 `json:"match_id"`
	QuestionID uint64 `json:"question_id"`
	PlayerID   uint64 `json:"player_id"`
	Answer     string `gorm:"type:varchar(5)" json:"answer"`
	IsCorrect  bool   `json:"is_correct"`
	Points     int    `json:"points"`
	Duration   int    `json:"duration"` // in seconds

	Question Question `gorm:"foreignKey:QuestionID" json:"question"`
}
