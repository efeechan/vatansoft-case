package controllers

import (
	"net/http"

	"github.com/efecan/vatansoft-case/config"
	"github.com/efecan/vatansoft-case/models"
	"github.com/gin-gonic/gin"
)

type AddressResponse struct {
	ID         uint   `json:"ID"`
	Street     string `json:"street"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type HospitalResponse struct {
	ID      uint            `json:"ID"`
	Name    string          `json:"name"`
	Phone   string          `json:"phone"`
	Address AddressResponse `json:"address"`
}

type DepartmentResponse struct {
	ID               uint             `json:"ID"`
	Name             string           `json:"name"`
	HospitalID       uint             `json:"hospital_id"`
	DepartmentTypeID uint             `json:"department_type_id"`
	DepartmentType   string           `json:"department_type"`
	Hospital         HospitalResponse `json:"hospital"`
}

type DoctorWithRelationsResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	HospitalID   uint   `json:"hospital_id"`
	DepartmentID uint   `json:"department_id"`
	Hospital     struct {
		ID      uint   `json:"id"`
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Address struct {
			ID         uint   `json:"id"`
			Street     string `json:"street"`
			City       string `json:"city"`
			PostalCode string `json:"postal_code"`
			Country    string `json:"country"`
		} `json:"address"`
	} `json:"hospital"`
	Department struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		HospitalID uint   `json:"hospital_id"`
	} `json:"department"`
}

// CreateDepartment godoc
// @Summary Create a department
// @Description Creates a new department under the admin's hospital
// @Tags Department
// @Accept json
// @Produce json
// @Param request body CreateDepartmentRequest true "Department creation payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /departments [post]
type CreateDepartmentRequest struct {
	DepartmentTypeID uint `json:"department_type_id" binding:"required"`
	HospitalID       uint `json:"hospital_id" binding:"required"`
}

func CreateDepartment(c *gin.Context) {
	var req struct {
		DepartmentTypeID uint `json:"department_type_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var admin models.User
	if err := config.DB.First(&admin, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var departmentType models.DepartmentType
	if err := config.DB.First(&departmentType, req.DepartmentTypeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department type not found"})
		return
	}

	department := models.Department{
		DepartmentTypeID: req.DepartmentTypeID,
		HospitalID:       admin.HospitalID,
	}

	if err := config.DB.Create(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create department"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department created successfully"})
}

// GetDepartments godoc
// @Summary List departments
// @Description Returns a list of all departments with their types and hospital info
// @Tags Department
// @Produce json
// @Success 200 {array} DepartmentResponse
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /departments [get]
func GetDepartments(c *gin.Context) {
	var departments []models.Department

	if err := config.DB.
		Preload("Hospital.Address").
		Preload("DepartmentType").
		Find(&departments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve departments"})
		return
	}

	var response []DepartmentResponse
	for _, dept := range departments {
		resp := DepartmentResponse{
			ID:               dept.ID,
			Name:             dept.Name,
			HospitalID:       dept.HospitalID,
			DepartmentTypeID: dept.DepartmentTypeID,
			DepartmentType:   dept.DepartmentType.Name,
			Hospital: HospitalResponse{
				ID:    dept.Hospital.ID,
				Name:  dept.Hospital.Name,
				Phone: dept.Hospital.Phone,
				Address: AddressResponse{
					ID:         dept.Hospital.Address.ID,
					Street:     dept.Hospital.Address.Street,
					City:       dept.Hospital.Address.City,
					PostalCode: dept.Hospital.Address.PostalCode,
					Country:    dept.Hospital.Address.Country,
				},
			},
		}
		response = append(response, resp)
	}

	c.JSON(http.StatusOK, response)
}

// GetDoctorsByDepartment godoc
// @Summary Get doctors by department ID
// @Description Returns all doctors in the specified department
// @Tags Department
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {array} DoctorWithRelationsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /departments/{id}/doctors [get]
func GetDoctorsByDepartment(c *gin.Context) {
	departmentID := c.Param("id")

	var doctors []models.Doctor
	err := config.DB.
		Preload("Hospital.Address").
		Preload("Department").
		Where("department_id = ?", departmentID).
		Find(&doctors).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}

	var response []DoctorWithRelationsResponse
	for _, d := range doctors {
		resp := DoctorWithRelationsResponse{
			ID:           d.ID,
			Name:         d.Name,
			Email:        d.Email,
			HospitalID:   d.HospitalID,
			DepartmentID: d.DepartmentID,
			Hospital: struct {
				ID      uint   `json:"id"`
				Name    string `json:"name"`
				Phone   string `json:"phone"`
				Address struct {
					ID         uint   `json:"id"`
					Street     string `json:"street"`
					City       string `json:"city"`
					PostalCode string `json:"postal_code"`
					Country    string `json:"country"`
				} `json:"address"`
			}{
				ID:    d.Hospital.ID,
				Name:  d.Hospital.Name,
				Phone: d.Hospital.Phone,
				Address: struct {
					ID         uint   `json:"id"`
					Street     string `json:"street"`
					City       string `json:"city"`
					PostalCode string `json:"postal_code"`
					Country    string `json:"country"`
				}{
					ID:         d.Hospital.Address.ID,
					Street:     d.Hospital.Address.Street,
					City:       d.Hospital.Address.City,
					PostalCode: d.Hospital.Address.PostalCode,
					Country:    d.Hospital.Address.Country,
				},
			},
			Department: struct {
				ID         uint   `json:"id"`
				Name       string `json:"name"`
				HospitalID uint   `json:"hospital_id"`
			}{
				ID:         d.Department.ID,
				Name:       d.Department.Name,
				HospitalID: d.Department.HospitalID,
			},
		}

		response = append(response, resp)
	}

	c.JSON(http.StatusOK, response)
}
