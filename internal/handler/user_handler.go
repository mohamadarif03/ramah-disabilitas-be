package handler

import (
	"net/http"
	"ramah-disabilitas-be/internal/service"
	"ramah-disabilitas-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

func UpdateAccessibility(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input service.AccessibilityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi input gagal.",
			"errors":  utils.FormatValidationError(err),
		})
		return
	}

	profile, err := service.UpdateAccessibilityProfile(userID.(uint64), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prepare user-friendly response messages (Step 3 Confirmation)
	var actions []string
	if profile.AISummary {
		actions = append(actions, "Materi akan otomatis diringkas (Mode Fokus)")
	}
	if profile.SubtitlesRequired {
		actions = append(actions, "Video akan selalu menampilkan subtitle")
	}
	if profile.ScreenReaderCompatible {
		actions = append(actions, "Fitur pembaca layar diaktifkan")
	}
	if profile.KeyboardNavigation {
		actions = append(actions, "Navigasi keyboard diaktifkan")
	}
	if profile.TextBasedSubmission {
		actions = append(actions, "Pengumpulan tugas via teks diizinkan")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Preferensi disabilitas berhasil disimpan",
		"data": gin.H{
			"profile":          profile,
			"confirmation_msg": "Terima kasih informasinya! Berdasarkan pilihanmu, kami telah menyiapkan:",
			"active_features":  actions,
		},
	})
}
