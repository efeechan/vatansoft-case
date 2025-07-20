package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/efecan/vatansoft-case/config"
	"github.com/efecan/vatansoft-case/models"
	"github.com/gin-gonic/gin"
)

type TitleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProfessionGroupResponse struct {
	ID     uint            `json:"id"`
	Name   string          `json:"name"`
	Titles []TitleResponse `json:"titles"`
}

// GetProfessionGroups godoc
// @Summary Get all profession groups with titles
// @Description Returns all profession groups and their associated titles, using Redis cache
// @Tags Profession Groups
// @Produce json
// @Success 200 {array} ProfessionGroupResponse
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /profession-groups [get]
func GetProfessionGroups(c *gin.Context) {
	const cacheKey = "profession_groups"

	cached, err := config.REDIS.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var response []ProfessionGroupResponse
		if json.Unmarshal([]byte(cached), &response) == nil {
			c.JSON(http.StatusOK, response)
			return
		}
	}

	var groups []models.ProfessionGroup
	if err := config.DB.Preload("Titles").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profession groups"})
		return
	}

	var response []ProfessionGroupResponse
	for _, g := range groups {
		var titles []TitleResponse
		for _, t := range g.Titles {
			titles = append(titles, TitleResponse{ID: t.ID, Name: t.Name})
		}
		response = append(response, ProfessionGroupResponse{
			ID:     g.ID,
			Name:   g.Name,
			Titles: titles,
		})
	}

	data, _ := json.Marshal(response)
	config.REDIS.Set(context.Background(), cacheKey, data, 24*time.Hour)

	c.JSON(http.StatusOK, response)
}
