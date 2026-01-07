package service

import (
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
)

type AccessibilityInput struct {
	// A=Vision, B=Hearing, C=Physical, D=Cognitive, E=Speech
	Categories []string `json:"categories" binding:"required"`
}

func UpdateAccessibilityProfile(userID uint64, input AccessibilityInput) (*model.AccessibilityProfile, error) {
	profile := &model.AccessibilityProfile{
		UserID: userID,
	}

	for _, category := range input.Categories {
		switch category {
		case "A", "vision", "penglihatan":
			profile.VisionImpaired = true
			profile.ScreenReaderCompatible = true
			profile.AudioDescription = true
		case "B", "hearing", "pendengaran":
			profile.HearingImpaired = true
			profile.SubtitlesRequired = true
			profile.VisualNotifications = true
		case "C", "physical", "motorik", "daksa":
			profile.PhysicalImpaired = true
			profile.KeyboardNavigation = true
			profile.VoiceCommand = true
		case "D", "cognitive", "fokus", "adhd", "disleksia":
			profile.CognitiveImpaired = true
			profile.AISummary = true
			profile.FocusMode = true
		case "E", "speech", "wicara", "bisu":
			profile.SpeechImpaired = true
			profile.TextBasedSubmission = true
		}
	}

	if err := repository.SaveAccessibilityProfile(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func GetAccessibilityProfile(userID uint64) (*model.AccessibilityProfile, error) {
	return repository.FindAccessibilityProfileByUserID(userID)
}
