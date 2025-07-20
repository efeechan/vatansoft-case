package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/efecan/vatansoft-case/config"
	"github.com/efecan/vatansoft-case/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type NewUserRequest struct {
	Name              string `json:"name" binding:"required"`
	Surname           string `json:"surname" binding:"required"`
	TCKN              string `json:"tckn" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Phone             string `json:"phone" binding:"required"`
	Password          string `json:"password" binding:"required"`
	Role              string `json:"role" binding:"required"`
	HospitalID        uint   `json:"hospital_id" binding:"required"`
	ProfessionGroupID uint   `json:"profession_group_id" binding:"required"`
	TitleID           uint   `json:"title_id" binding:"required"`
}

type UserResponse struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	Surname           string `json:"surname"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	TCKN              string `json:"tckn"`
	Role              string `json:"role"`
	HospitalID        uint   `json:"hospital_id"`
	ProfessionGroupID uint   `json:"profession_group_id"`
	ProfessionGroup   string `json:"profession_group"`
	TitleID           uint   `json:"title_id"`
	Title             string `json:"title"`
}

// CreateUser godoc
// @Summary Create a new user (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param user body NewUserRequest true "New user data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var req NewUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:              req.Name,
		Surname:           req.Surname,
		TCKN:              req.TCKN,
		Email:             req.Email,
		Phone:             req.Phone,
		Password:          string(hashedPassword),
		Role:              strings.ToLower(req.Role),
		HospitalID:        req.HospitalID,
		ProfessionGroupID: req.ProfessionGroupID,
		TitleID:           req.TitleID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// GetUsers godoc
// @Summary Get all users in the current user's hospital
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users [get]
func GetUsers(c *gin.Context) {
	hospitalID := c.GetInt("hospitalID")

	var users []models.User
	if err := config.DB.Where("hospital_id = ?", hospitalID).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// ListUsers godoc
// @Summary List users with filtering and pagination (admin only)
// @Tags Users
// @Produce json
// @Param page query int false "Page number"
// @Param name query string false "Filter by name"
// @Param surname query string false "Filter by surname"
// @Param tckn query string false "Filter by TCKN"
// @Param profession_group_id query string false "Filter by profession group ID"
// @Param title_id query string false "Filter by title ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /listusers [get]
func ListUsers(c *gin.Context) {
	hospitalID := c.GetInt("hospitalID")

	pageStr := c.DefaultQuery("page", "1")
	name := c.Query("name")
	surname := c.Query("surname")
	tckn := c.Query("tckn")
	professionGroupID := c.Query("profession_group_id")
	titleID := c.Query("title_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	var users []models.User
	query := config.DB.Preload("ProfessionGroup").Preload("Title").Where("hospital_id = ?", hospitalID)

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if surname != "" {
		query = query.Where("surname ILIKE ?", "%"+surname+"%")
	}
	if tckn != "" {
		query = query.Where("tckn = ?", tckn)
	}
	if professionGroupID != "" {
		query = query.Where("profession_group_id = ?", professionGroupID)
	}
	if titleID != "" {
		query = query.Where("title_id = ?", titleID)
	}

	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var response []UserResponse
	for _, u := range users {
		response = append(response, UserResponse{
			ID:                u.ID,
			Name:              u.Name,
			Surname:           u.Surname,
			Email:             u.Email,
			Phone:             u.Phone,
			TCKN:              u.TCKN,
			Role:              u.Role,
			HospitalID:        u.HospitalID,
			ProfessionGroupID: u.ProfessionGroupID,
			ProfessionGroup:   u.ProfessionGroup.Name,
			TitleID:           u.TitleID,
			Title:             u.Title.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"data":  response,
	})
}

// UpdateUser godoc
// @Summary Update a user's info (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body map[string]string true "Updated user data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		TCKN     string `json:"tckn"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user.Name = input.Name
	user.Surname = input.Surname
	user.Email = input.Email
	user.Phone = input.Phone
	user.TCKN = input.TCKN
	user.Role = input.Role

	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser godoc
// @Summary Delete a user (admin only)
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
