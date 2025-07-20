package middlewares

import (
	"net/http"
	"strings"

	"github.com/efecan/vatansoft-case/config"
	"github.com/efecan/vatansoft-case/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing auth header"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth format"})
		return
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET", "devsecret")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
		return
	}

	name, _ := claims["name"].(string)
	role, _ := claims["role"].(string)
	hospitalID, _ := claims["hospital_id"].(float64)

	c.Set("userID", int(sub))
	c.Set("userName", name)
	c.Set("userRole", role)
	c.Set("hospitalID", int(hospitalID))

	c.Next()
}

func RequireAdmin(c *gin.Context) {
	role, exists := c.Get("userRole")
	if !exists || role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}
	c.Next()
}

func IsAdminOfSameHospitalOrAbort(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	userHospitalID := c.GetInt("hospitalID")

	id := c.Param("id")

	var target models.User
	if err := config.DB.First(&target, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if userRole != "admin" || target.HospitalID != uint(userHospitalID) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
}
