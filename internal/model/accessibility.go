package model

type AccessibilityProfile struct {
	UserID uint64 `gorm:"primaryKey" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// User Conditions ( Needs )
	VisionImpaired    bool `json:"vision_impaired"`
	HearingImpaired   bool `json:"hearing_impaired"`
	PhysicalImpaired  bool `json:"physical_impaired"`
	CognitiveImpaired bool `json:"cognitive_impaired"`
	SpeechImpaired    bool `json:"speech_impaired"`

	// System Features ( Actions )
	ScreenReaderCompatible bool `json:"screen_reader_compatible"`
	AudioDescription       bool `json:"audio_description"`
	SubtitlesRequired      bool `json:"subtitles_required"`
	VisualNotifications    bool `json:"visual_notifications"`
	KeyboardNavigation     bool `json:"keyboard_navigation"`
	VoiceCommand           bool `json:"voice_command"`
	AISummary              bool `json:"ai_summary"`
	FocusMode              bool `json:"focus_mode"`
	TextBasedSubmission    bool `json:"text_based_submission"`
}
