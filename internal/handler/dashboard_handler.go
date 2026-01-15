package handler

import (
	"net/http"
	"ramah-disabilitas-be/internal/service"

	"github.com/gin-gonic/gin"
)

func GetLecturerDashboardStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Assuming the user is a lecturer. Middleware might handle role check, or we can check here.
	// But usually if they hit this endpoint they are authenticated.
	// We can trust the userID.

	stats, err := service.GetLecturerDashboardStats(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dashboard statistics retrieved successfully",
		"data":    stats,
	})
}
