package model

import (
	"time"
)

type MaterialCompletion struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64    `gorm:"uniqueIndex:idx_user_material" json:"user_id"`
	MaterialID uint64    `gorm:"uniqueIndex:idx_user_material" json:"material_id"`
	Completed  bool      `json:"completed"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
