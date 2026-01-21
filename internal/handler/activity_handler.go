package handler

import (
	"net/http"
	"ramah-disabilitas-be/internal/service"
	"ramah-disabilitas-be/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLecturerActivities(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		utils.JSONResponse(c, http.StatusUnauthorized, "Unauthorized", "error", nil)
		return
	}
	userID := userIDStr.(uint64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	resp, err := service.GetLecturerActivities(userID, page, limit)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "Failed to fetch activities", "error", nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "Success fetch activities", "success", resp)
}
