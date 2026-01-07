package middleware

import (
	"net/http"
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func LecturerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Header otorisasi diperlukan"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")
		if len(tokenString) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid"})
			c.Abort()
			return
		}

		token, err := utils.ValidateToken(tokenString[1])
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kadaluarsa"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Klaim token tidak valid"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || (role != string(model.RoleLecturer) && role != string(model.RoleAdmin)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Hanya dosen yang diizinkan"})
			c.Abort()
			return
		}

		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("userID", uint64(userID))
		}
		c.Set("role", role)

		c.Next()
	}
}
