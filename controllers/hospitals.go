package controllers

import (
	"net/http"

	"github.com/efecan/vatansoft-case/config"
	"github.com/efecan/vatansoft-case/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type HospitalRegisterRequest struct {
	HospitalName string `json:"hospital_name"`
	Phone        string `json:"phone"`
	Address      struct {
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postal_code"`
		Country    string `json:"country"`
	} `json:"address"`
	AdminUser struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"admin_user"`
}

// HospitalRegister godoc
// @Summary Register a new hospital and admin user
// @Description Creates a new hospital with address and registers an admin user for it
// @Tags hospitals
// @Accept json
// @Produce json
// @Param request body controllers.HospitalRegisterRequest true "Hospital and Admin registration info"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hospitals/register [post]
func HospitalRegister(c *gin.Context) {
	var req HospitalRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	db := config.DB

	address := models.Address{
		Street:     req.Address.Street,
		City:       req.Address.City,
		PostalCode: req.Address.PostalCode,
		Country:    req.Address.Country,
	}
	if err := db.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save address"})
		return
	}

	hospital := models.Hospital{
		Name:      req.HospitalName,
		Phone:     req.Phone,
		AddressID: address.ID,
	}
	if err := db.Create(&hospital).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save hospital"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.AdminUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	admin := models.User{
		Name:       req.AdminUser.Name,
		Email:      req.AdminUser.Email,
		Password:   string(hash),
		HospitalID: hospital.ID,
	}
	if err := db.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save admin user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "hospital & admin created"})

}

// GetHospitals godoc
// @Summary List all hospitals with their address and admin users
// @Description Returns a list of all registered hospitals including address and admin info
// @Tags hospitals
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /hospitals [get]
func GetHospitals(c *gin.Context) {
	var hospitals []models.Hospital

	err := config.DB.Preload("Address").Preload("Users").Find(&hospitals).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hospitals"})
		return
	}

	var response []gin.H
	for _, h := range hospitals {
		var admins []gin.H
		for _, u := range h.Users {
			admins = append(admins, gin.H{
				"id":    u.ID,
				"name":  u.Name,
				"email": u.Email,
			})
		}

		response = append(response, gin.H{
			"id":      h.ID,
			"name":    h.Name,
			"phone":   h.Phone,
			"address": h.Address,
			"admins":  admins,
		})
	}

	c.JSON(http.StatusOK, response)
}
