package repository

import (
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/pkg/database"

	"gorm.io/gorm/clause"
)

func SaveAccessibilityProfile(profile *model.AccessibilityProfile) error {
	// Upsert: On conflict (user_id), update everything
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(profile).Error
}

func FindAccessibilityProfileByUserID(userID uint64) (*model.AccessibilityProfile, error) {
	var profile model.AccessibilityProfile
	err := database.DB.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}
