package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/efecan/vatansoft-case/config"
	"github.com/efecan/vatansoft-case/models"
	"github.com/gin-gonic/gin"
)

type DistrictResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CityResponse struct {
	ID        uint               `json:"id"`
	Name      string             `json:"name"`
	Districts []DistrictResponse `json:"districts"`
}

// GetCities godoc
// @Summary List cities with districts
// @Description Returns a list of all cities and their districts, with Redis caching
// @Tags Location
// @Produce json
// @Success 200 {array} CityResponse
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /cities [get]
func GetCities(c *gin.Context) {
	cacheKey := "cities_with_districts"

	cached, err := config.REDIS.Get(c, cacheKey).Result()
	if err == nil {
		var response []CityResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			c.JSON(http.StatusOK, response)
			return
		}
	}

	var cities []models.City
	if err := config.DB.Preload("Districts").Find(&cities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cities"})
		return
	}

	var response []CityResponse
	for _, city := range cities {
		var districts []DistrictResponse
		for _, d := range city.Districts {
			districts = append(districts, DistrictResponse{
				ID:   d.ID,
				Name: d.Name,
			})
		}

		response = append(response, CityResponse{
			ID:        city.ID,
			Name:      city.Name,
			Districts: districts,
		})
	}

	data, _ := json.Marshal(response)
	config.REDIS.Set(c, cacheKey, data, 24*time.Hour)

	c.JSON(http.StatusOK, response)
}
