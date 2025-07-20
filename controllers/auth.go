package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"time"

	"github.com/efecan/vatansoft-case/config"
	"github.com/efecan/vatansoft-case/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Registers a new hospital and its admin user
// @Description Creates a hospital and registers an admin for it. Only one admin per hospital is allowed.
// @Tags Auth
// @Accept json
// @Produce json
// @Param registration body object true "Hospital and admin registration info"
// @Success 200 {object} map[string]string "Hospital and admin user registered successfully"
// @Failure 400 {object} map[string]string "Invalid input or hospital already exists"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /register [post]
func Register(c *gin.Context) {
	var req struct {
		Hospital struct {
			Name      string `json:"name" binding:"required"`
			TaxNumber string `json:"tax_number" binding:"required"`
			Email     string `json:"email" binding:"required,email"`
			Phone     string `json:"phone" binding:"required"`
			Address   struct {
				ProvinceID uint   `json:"province_id" binding:"required"`
				DistrictID uint   `json:"district_id" binding:"required"`
				Street     string `json:"street" binding:"required"`
			} `json:"address"`
		} `json:"hospital"`
		Admin struct {
			Name              string `json:"name" binding:"required"`
			Surname           string `json:"surname" binding:"required"`
			TCKN              string `json:"tc_no" binding:"required"`
			Email             string `json:"email" binding:"required,email"`
			Phone             string `json:"phone" binding:"required"`
			Password          string `json:"password" binding:"required"`
			Role              string `json:"role" binding:"required"` // must be "admin"
			ProfessionGroupID uint   `json:"profession_group_id" binding:"required"`
			TitleID           uint   `json:"title_id" binding:"required"`
		} `json:"admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if strings.ToLower(req.Admin.Role) != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only admin registration is allowed on this endpoint"})
		return
	}

	var existingHospital models.Hospital
	if err := config.DB.Where("tax_number = ? OR email = ? OR phone = ?", req.Hospital.TaxNumber, req.Hospital.Email, req.Hospital.Phone).
		First(&existingHospital).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hospital already exists with provided tax/email/phone"})
		return
	}

	address := models.Address{
		ProvinceID: req.Hospital.Address.ProvinceID,
		DistrictID: req.Hospital.Address.DistrictID,
		Street:     req.Hospital.Address.Street,
	}
	if err := config.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	hospital := models.Hospital{
		Name:      req.Hospital.Name,
		TaxNumber: req.Hospital.TaxNumber,
		Email:     req.Hospital.Email,
		Phone:     req.Hospital.Phone,
		AddressID: address.ID,
	}
	if err := config.DB.Create(&hospital).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hospital"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var existingAdmin models.User
	if err := config.DB.Where("hospital_id = ? AND role = ?", hospital.ID, "admin").First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This hospital already has an admin user"})
		return
	}

	admin := models.User{
		Name:              req.Admin.Name,
		Surname:           req.Admin.Surname,
		TCKN:              req.Admin.TCKN,
		Email:             req.Admin.Email,
		Phone:             req.Admin.Phone,
		Password:          string(hashedPassword),
		Role:              "admin",
		HospitalID:        hospital.ID,
		ProfessionGroupID: req.Admin.ProfessionGroupID,
		TitleID:           req.Admin.TitleID,
	}
	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hospital and admin user registered successfully"})
}

// Login godoc
// @Summary Logs in a user and returns a JWT token
// @Description Authenticates the user by email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body object true "User credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "invalid credentials"
// @Router /login [post]
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         user.ID,
		"name":        user.Name,
		"role":        user.Role,
		"hospital_id": user.HospitalID,
		"exp":         time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET", "devsecret")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func generateResetCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// RequestPasswordReset godoc
// @Summary Request password reset code
// @Description Sends a reset code for the given phone number. The code is returned in the response (simulating SMS).
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body object true "Phone number for password reset"
// @Success 200 {object} map[string]interface{} "Reset code generated"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Error creating reset code"
// @Router /auth/request-password-reset [post]
func RequestPasswordReset(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required and must be valid"})
		return
	}

	var user models.User
	if err := config.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found for the provided phone number"})
		return
	}

	code := generateResetCode()
	key := "reset_code:" + req.Phone
	expiration := 10 * time.Minute

	err := config.REDIS.Set(context.Background(), key, code, expiration).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating reset code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reset code has been sent to your phone number. Please use this code to reset your password.",
		"code":    code,
	})
}

// ResetPassword godoc
// @Summary Reset password using verification code
// @Description Resets the user's password after verifying the code sent to their phone
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body object true "Reset password request"
// @Success 200 {object} map[string]string "Password reset successfully"
// @Failure 400 {object} map[string]string "Invalid request or code mismatch"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Failed to update password"
// @Router /auth/reset-password [post]
func ResetPassword(c *gin.Context) {
	var req struct {
		Phone       string `json:"phone" binding:"required"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
		ConfirmPass string `json:"confirm_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if req.NewPassword != req.ConfirmPass {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	cacheKey := "reset_code:" + req.Phone
	storedCode, err := config.REDIS.Get(c, cacheKey).Result()
	if err != nil || storedCode != req.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired code"})
		return
	}

	var user models.User
	if err := config.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password", "details": err.Error()})
		return
	}

	config.REDIS.Del(c, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
